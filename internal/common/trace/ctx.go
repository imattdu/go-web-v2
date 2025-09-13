package trace

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
)

var traceCtxKey = "x-trace"

func SetTrace(ctx context.Context, trace *Trace) {
	cctx.Set(ctx, traceCtxKey, trace)
}

func GetTrace(ctx context.Context) *Trace {
	v, _ := cctx.Get(ctx, traceCtxKey)
	if val, ok := v.(*Trace); ok {
		return val
	}
	return &Trace{}
}
