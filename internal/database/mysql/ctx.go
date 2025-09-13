package mysql

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"time"
)

var callStatsKey = "mysqlCallStatsKey"

type CallStats struct {
	Params interface{}
	Start  time.Time
}

func SetCallStats(ctx context.Context, stats *CallStats) {
	cctx.Set(ctx, callStatsKey, stats)
}

func SetCallStatsClone(ctx context.Context, stats *CallStats) context.Context {
	ctx = cctx.CloneWithoutDeadline(ctx)
	cctx.Set(ctx, callStatsKey, stats)
	return ctx
}

func GetCallStats(ctx context.Context) *CallStats {
	v, _ := cctx.Get(ctx, callStatsKey)
	if val, ok := v.(*CallStats); ok {
		return val
	}
	return &CallStats{}
}
