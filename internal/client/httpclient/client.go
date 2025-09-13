package httpclient

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/logger"
	"github.com/imattdu/go-web-v2/internal/common/trace"

	"github.com/go-resty/resty/v2"
)

var GlobalClient *Client

func Init() {
	GlobalClient = NewClient()
}

// NewClient 创建一个新的 httpclient
func NewClient() *Client {
	transport := &http.Transport{
		// 最大空闲连接数
		MaxIdleConns: 100,
		// 每个主机最大空闲连接数
		MaxIdleConnsPerHost: 10,
		// 空闲连接存活时间
		IdleConnTimeout: 90 * time.Second,
		// TLS 握手超时
		TLSHandshakeTimeout: 10 * time.Second,
		// 每个请求等待连接可用的时间
		ExpectContinueTimeout: 1 * time.Second,
		// 自定义拨号器（可以设置连接超时）
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // 连接超时
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	c := resty.New().
		SetTransport(transport).
		SetRetryCount(0).
		OnError(func(request *resty.Request, err error) {
			if err == nil {
				return
			}
			var mErr errorx.MErr
			if errors.As(err, &mErr) {
				return
			}

			var (
				ctx    = request.Context()
				req    = GetHttpRequest(ctx)
				logMap = map[string]interface{}{
					logger.KURL:         request.URL,
					logger.KHeaders:     request.Header,
					logger.KRequestBody: request.Body,

					logger.KProcTime: time.Now().Sub(request.Time) / time.Millisecond,
					logger.KAttempt:  req.Stats.attempt,
					logger.KRetries:  req.Retries,
				}
			)
			err = errorx.New(errorx.ErrOptions{
				ErrMeta: errorx.ErrMeta{
					ServiceType: errorx.ServiceTypeService,
					Service:     req.Service,
					ErrType:     errorx.ErrTypeSys,
				},
				Err: err,
			})
			req.Stats.lastError = err

			rpcFinal := req.Stats.rpcFinal
			if !rpcFinal {
				rpcFinal = !shouldRetry(req, nil, err)
			}
			logMap[logger.KRPCFinal] = rpcFinal
			logMap[logger.KErr] = err
			logMap[logger.KErrMsg] = err.Error()
			logger.Warn(ctx, logger.TagHttpFailure, logMap)
		}).
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			var (
				ctx = r.Context()
				req = GetHttpRequest(ctx)
			)
			req.Stats.rpcFinal = req.Stats.attempt == req.Retries
			return nil
		}).
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			var (
				ctx    = r.Request.Context()
				req    = GetHttpRequest(ctx)
				logMap = map[string]interface{}{
					logger.KURL:          r.Request.URL,
					logger.KHeaders:      r.Request.Header,
					logger.KRequestBody:  r.Request.Body,
					logger.KResponseBody: string(r.Body()),

					logger.KProcTime: r.Time() / time.Millisecond,
					logger.KAttempt:  req.Stats.attempt,
					logger.KRetries:  req.Retries,
				}
				err error
			)
			defer func() {
				req.Stats.lastError = err
				rpcFinal := req.Stats.rpcFinal
				if !rpcFinal {
					rpcFinal = !shouldRetry(req, r, err)
				}
				logMap[logger.KRPCFinal] = rpcFinal
				mErr := errorx.Get(err, false)
				if mErr != nil {
					logMap[logger.KErrMsg] = mErr.FinalMsg
				}
				if mErr != nil && mErr.ErrType == errorx.ErrTypeSys {
					logger.Warn(ctx, logger.TagHttpFailure, logMap)
				} else {
					logger.Info(ctx, logger.TagHttpSuccess, logMap)
				}
			}()

			if r.StatusCode() != 200 {
				err = errorx.New(errorx.ErrOptions{
					ErrMeta: errorx.ErrMeta{
						ServiceType: errorx.ServiceTypeService,
						Service:     req.Service,
					},
					Err:  errors.New(r.Status()),
					Code: r.StatusCode(),
				})
				return err
			}
			if req.IsError != nil {
				err = errorx.New(errorx.ErrOptions{
					ErrMeta: errorx.ErrMeta{
						ServiceType: errorx.ServiceTypeService,
						Service:     req.Service,
					},
					Err: req.IsError(r.RawResponse),
				})
				return err
			}
			return nil
		})
	return &Client{client: c}
}

func shouldRetry(req *HttpRequest, response *resty.Response, err error) bool {
	if err == nil {
		return false
	}
	if response == nil {
		return true
	}
	mErr := errorx.Get(err, false)
	if mErr.ErrType == errorx.ErrTypeSys {
		return true
	}
	if req.RetryIf != nil {
		return req.RetryIf(response.RawResponse, err)
	}
	return false
}

func do(ctx context.Context, request *HttpRequest) error {
	newCtx := cctx.CloneWithoutDeadline(ctx)
	SetHttpRequest(newCtx, request)
	for i := 0; i <= request.Retries; i++ {
		request.Stats.attempt = i
		request.Stats.rpcFinal = i == request.Retries

		// 每次循环都用独立作用域，确保 cancel 在本次结束时调用
		err := func() error {
			if request.Timeout > 0 {
				newCtx = cctx.CloneWithoutDeadline(ctx)
				var cancel context.CancelFunc
				newCtx, cancel = context.WithTimeout(newCtx, request.Timeout)
				defer cancel()
			}

			t := trace.GetTrace(newCtx)
			t = t.Copy()
			t.UpdateParentSpanID()
			trace.SetTrace(newCtx, t)

			var (
				err     error
				headers = trace.NewHeader(newCtx, t)
			)
			for k, v := range request.Headers {
				headers[k] = v
			}
			r := GlobalClient.client.R().
				SetContext(newCtx).
				SetHeaders(headers).
				SetBody(request.JSONBody).
				SetResult(request.ResponseBody)
			switch request.method {
			case http.MethodGet:
				_, _ = r.Get(request.URL)
			case http.MethodPost:
				_, _ = r.Post(request.URL)
			}
			err = request.Stats.lastError

			if request.Stats.rpcFinal {
				time.Sleep(time.Duration(i+1) * 100 * time.Millisecond) // 退避
			}
			return err
		}()
		// 没有重试就返回
		if request.Stats.rpcFinal {
			return err
		}
	}
	return nil
}
