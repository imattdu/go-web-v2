package redis

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"github.com/imattdu/go-web-v2/internal/common/trace"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/errorx"
)

type retryConf struct {
	retries int
	delay   time.Duration
}

type retryClient struct {
	retryConf
}

func newRetryClient(conf retryConf) *retryClient {
	return &retryClient{
		retryConf: conf,
	}
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	mErr := errorx.Get(err, false)
	return mErr.ErrType == errorx.ErrTypeSys
}

func (r *retryClient) RetryCmd(ctx context.Context, CmdFunc func(_ context.Context) error) error {
	var err error
	for i := 0; i <= r.retries; i++ {
		c := cctx.With(ctx, cctx.RedisCallStatsKey, &CallStats{})

		t := cctx.GetOrNewAs[*trace.Trace](c, cctx.TraceKey, func() *trace.Trace {
			return trace.New(nil)
		})
		t = t.Copy()
		t.UpdateParentSpanID()
		c = cctx.With(c, cctx.TraceKey, t)

		err = CmdFunc(c)
		if !shouldRetry(err) {
			return err
		}
		time.Sleep(time.Duration(i+1) * r.delay) // 退避
	}
	return err
}
