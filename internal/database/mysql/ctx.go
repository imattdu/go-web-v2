package mysql

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"time"
)

type ctxKey string

var callStatsKey ctxKey = "mysqlCallStatsKey"

type CallStats struct {
	Params interface{}
	Start  time.Time
}

func WithCallStatsCtx(ctx context.Context, stats CallStats) context.Context {
	ctx = cctx.Get(ctx)
	return context.WithValue(ctx, callStatsKey, stats)
}

func CallStatsFromCtx(ctx context.Context) CallStats {
	if v, ok := ctx.Value(callStatsKey).(CallStats); ok {
		return v
	}
	return CallStats{}
}
