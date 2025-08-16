package httpclientresty

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
		AddRetryCondition(func(r *resty.Response, err error) bool {
			var (
				ctx = r.Request.Context()
				req = ReqFromCtx(ctx)
			)
			return shouldRetry(req, r, err)
		}).
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
				req    = ReqFromCtx(ctx)
				logMap = map[string]interface{}{
					logger.KURL:          request.URL,
					logger.KHeaders:      request.Header,
					logger.KRequestBody:  request.Body,
					logger.KResponseBody: request.Body,
					logger.KProcTime:     time.Now().Sub(request.Time) / time.Millisecond,
					logger.KRetry:        req.Stats.retry,
					logger.KRetryCount:   req.Meta.RetryCount,
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
			isRpcFinal := req.Stats.isRpcFinal
			if !isRpcFinal {
				isRpcFinal = !shouldRetry(req, nil, err)
			}
			logMap[logger.KIsRPCFinal] = isRpcFinal
			logMap[logger.KErr] = err
			logMap[logger.KErrMsg] = err.Error()
			logger.Warn(ctx, logger.TagHttpFailure, logMap)
		}).
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			var (
				ctx = r.Context()
				req = ReqFromCtx(ctx)
			)
			req.Stats.isRpcFinal = req.Stats.retry == req.Meta.RetryCount
			ctx = WithReqCtx(ctx, req)
			r.SetContext(ctx)
			return nil
		}).
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			var (
				ctx    = r.Request.Context()
				req    = ReqFromCtx(ctx)
				logMap = map[string]interface{}{
					logger.KURL:          r.Request.URL,
					logger.KHeaders:      r.Request.Header,
					logger.KRequestBody:  r.Request.Body,
					logger.KResponseBody: r.Body(),
					logger.KProcTime:     r.Time() / time.Millisecond,
					logger.KRetry:        req.Stats.retry,
					logger.KRetryCount:   req.Meta.RetryCount,
				}
				err error
			)
			defer func() {
				isRpcFinal := req.Stats.isRpcFinal
				if !isRpcFinal {
					isRpcFinal = !shouldRetry(req, r, err)
				}
				logMap[logger.KIsRPCFinal] = isRpcFinal
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
			if req.Meta.IsError != nil {
				err = errorx.New(errorx.ErrOptions{
					ErrMeta: errorx.ErrMeta{
						ServiceType: errorx.ServiceTypeService,
						Service:     req.Service,
					},
					Err: req.Meta.IsError(r.RawResponse),
				})
				return err
			}
			return nil
		})
	return &Client{client: c}
}

func shouldRetry(req *Req, response *resty.Response, err error) bool {
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
	if req.Meta.RetryIf != nil {
		return req.Meta.RetryIf(response.RawResponse, err)
	}
	return false
}

func do(ctx context.Context, params *Req) error {
	for i := 0; i <= params.Meta.RetryCount; i++ {
		params.Stats.retry = i
		var (
			newCtx  = WithReqCtx(ctx, params)
			isRetry bool
		)
		// 每次循环都用独立作用域，确保 cancel 在本次结束时调用
		err := func() error {
			if params.Meta.Timeout > 0 {
				var cancel context.CancelFunc
				newCtx, cancel = context.WithTimeout(newCtx, params.Meta.Timeout)
				defer cancel()
			}

			t := cctx.TraceFromCtxOrNew(ctx, nil)
			t = t.Copy()
			t.UpdateParentSpanID()
			newCtx = cctx.WithTraceCtx(newCtx, t)

			var (
				response *resty.Response
				err      error
				headers  = trace.NewHeader(newCtx, t)
			)
			request := GlobalClient.client.R().
				SetContext(newCtx).
				SetHeaders(headers).
				SetBody(params.Meta.RequestBody).
				SetResult(params.Meta.ResponseBody)
			switch params.Meta.method {
			case http.MethodGet:
				response, err = request.Get(params.Meta.URL)
			case http.MethodPost:
				response, err = request.Post(params.Meta.URL)
			}

			isRetry = shouldRetry(params, response, err)
			if i == params.Meta.RetryCount {
				isRetry = false
			}
			params.Stats.isRpcFinal = !isRetry
			if isRetry {
				time.Sleep(time.Duration(i+1) * 100 * time.Millisecond) // 退避
			}
			return err
		}()
		// 没有重试就返回
		if !isRetry {
			return err
		}
	}
	return nil
}
