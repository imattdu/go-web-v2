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
	var (
		err error
		c   = cctx.CloneWithoutDeadline(ctx)
	)
	for i := 0; i <= r.retries; i++ {
		SetCallStats(c, CallStats{
			Attempt: i,
			Retries: r.retries,
		})

		t := trace.GetTrace(c)
		t = t.Copy()
		t.UpdateParentSpanID()
		trace.SetTrace(c, t)
		err = CmdFunc(c)
		if !shouldRetry(err) {
			return err
		}
		time.Sleep(time.Duration(i+1) * r.delay) // 退避
	}
	return err
}
