package httpclient

import (
	"context"
	"errors"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	errorx2 "github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/logger"
	"github.com/imattdu/go-web-v2/internal/common/trace"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

func do(ctx context.Context, req *Req) error {
	var err error
	for attempt := 0; attempt <= req.Meta.RetryCount; attempt++ {
		req.Stats.retry = attempt
		t := cctx.TraceFromCtxOrNew(ctx, nil).Copy()
		t.UpdateParentSpanID()
		ctx = cctx.WithTraceCtx(ctx, t)

		if err = prepareRequest(ctx, req); err != nil {
			break
		}
		req.Stats.rawResponse, req.Stats.responseText, req.Stats.errs = req.client.EndStruct(&req.Meta.ResponseBody)
		req.Stats.duration = time.Since(req.Stats.startTime)

		err = validateResponse(req)
		isRetry := req.shouldRetry(err) && attempt != req.Meta.RetryCount
		req.Stats.isRpcFinal = !isRetry
		collect(ctx, req, err)
		if !isRetry {
			break
		}
		time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond) // 递增退避
	}
	return err
}

func prepareRequest(ctx context.Context, req *Req) error {
	if req.Meta.Headers == nil {
		req.Meta.Headers = make(map[string]string, 16)
	}
	t := cctx.TraceFromCtxOrNew(ctx, nil)
	traceHeaders := trace.NewHeader(ctx, t)
	for k, v := range traceHeaders {
		req.Meta.Headers[k] = v
	}
	req.client = gorequest.New()
	if req.Meta.Timeout > 0 {
		req.client.Timeout(req.Meta.Timeout)
	}
	switch strings.ToUpper(req.Meta.Method) {
	case http.MethodPost:
		req.client = req.client.Post(req.Meta.URL).Send(req.Meta.RequestBody)
	default:
		req.client = req.client.Get(req.Meta.URL).Send(req.Meta.RequestBody)
	}
	for k, v := range req.Meta.Headers {
		req.client = req.client.Set(k, v)
	}
	req.Stats.startTime = time.Now()
	return nil
}

func validateResponse(req *Req) error {
	req.Stats.code = errorx2.Success.Code
	if len(req.Stats.errs) > 0 {
		return errorx2.New(errorx2.NewQuery{
			ErrMeta: errorx2.ErrMeta{
				ServiceType: errorx2.ServiceTypeService,
				Service:     req.Service,
				ErrType:     errorx2.ErrTypeSys,
			},
			Err: errors.New(errorx2.Errs2Msg(req.Stats.errs)),
		})
	}
	if req.Stats.rawResponse != nil && req.Stats.rawResponse.StatusCode != http.StatusOK {
		return errorx2.New(errorx2.NewQuery{
			ErrMeta: errorx2.ErrMeta{
				ServiceType: errorx2.ServiceTypeService,
				Service:     req.Service,
				ErrType:     errorx2.ErrTypeSys,
			},
			Err: errors.New(req.Stats.rawResponse.Status),
		})
	}
	if req.Meta.OnError != nil {
		if err := req.Meta.OnError(req.Stats.rawResponse, req.Stats.responseText); err != nil {
			return err
		}
	}
	return nil
}

func (r Req) shouldRetry(err error) bool {
	if r.Meta.RetryIf != nil {
		if ok := r.Meta.RetryIf(r.Stats.rawResponse, r.Stats.responseText, err); ok {
			return true
		}
	}
	// 非 200需要重试
	if r.Stats.rawResponse != nil && r.Stats.rawResponse.StatusCode != http.StatusOK {
		return true
	}
	mErr := errorx2.Get(err, false)
	if mErr != nil && mErr.ErrType == errorx2.ErrTypeSys {
		return true
	}
	return false
}

func collect(ctx context.Context, req *Req, err error) {
	var (
		logMap = map[string]interface{}{
			logger.KURL:          req.Meta.URL,
			logger.KHeaders:      req.Meta.Headers,
			logger.KRequestBody:  req.Meta.RequestBody,
			logger.KResponseBody: req.Meta.ResponseBody,
			logger.KResponseText: string(req.Stats.responseText),

			logger.KProcTime:   req.Stats.duration.Milliseconds(),
			logger.KCode:       req.Stats.code,
			logger.KIsRPCFinal: req.Stats.isRpcFinal,
			logger.KRetry:      req.Stats.retry,
			logger.KRetryCount: req.Meta.RetryCount,
		}
		mErr = errorx2.Get(err, false)
	)
	if mErr != nil {
		logMap[logger.KErr] = mErr.FinalMsg
	}

	if mErr != nil && mErr.ErrType == errorx2.ErrTypeSys {
		logger.Warn(ctx, logger.TagHttpFailure, logMap)
	} else {
		logger.Info(ctx, logger.TagHttpSuccess, logMap)
	}
}
