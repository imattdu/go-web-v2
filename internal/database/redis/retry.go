package redis

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"time"
)

type RetryClient struct {
	retries int
	delay   time.Duration
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	mErr := errorx.Get(err, false)
	return mErr.ErrType == errorx.ErrTypeSys
}

func (r *RetryClient) RetryCmd(ctx context.Context, CmdFunc func(_ context.Context) error) error {
	var err error
	for i := 0; i <= r.retries; i++ {
		newCtx := WithCallStatsCtx(ctx, CallStats{
			Attempt: i,
			Retries: r.retries,
		})

		trace := cctx.TraceFromCtxOrNew(newCtx, nil)
		trace = trace.Copy()
		trace.UpdateParentSpanID()
		newCtx = cctx.WithTraceCtx(newCtx, trace)
		err = CmdFunc(newCtx)
		if !shouldRetry(err) {
			return err
		}
		time.Sleep(time.Duration(i+1) * r.delay) // 退避
	}
	return err
}
