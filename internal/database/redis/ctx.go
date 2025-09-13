package redis

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"time"
)

type CallStats struct {
	Attempt int
	Retries int
	Start   time.Time
}

var callStatsCtxKey = "redisCallStats"

func SetCallStats(ctx context.Context, stats CallStats) {
	cctx.Set(ctx, callStatsCtxKey, stats)
}

func GetCallStats(ctx context.Context) CallStats {
	v, _ := cctx.Get(ctx, callStatsCtxKey)
	val, ok := v.(CallStats)
	if ok {
		return val
	}
	return CallStats{}
}
