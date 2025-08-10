package cctx

import (
	"context"

	"github.com/gin-gonic/gin"
	tracex "github.com/imattdu/go-web-v2/internal/common/trace"
	"net/http"
)

type ctxKey string

var (
	traceCtxKey ctxKey = "traceCtxKey"
	ginCtxKey   ctxKey = "ginCtxKey"
)

func WithTraceCtx(ctx context.Context, trace *tracex.Trace) context.Context {
	ctx = Get(ctx)
	return context.WithValue(ctx, traceCtxKey, trace)
}

func TraceFromCtxOrNew(ctx context.Context, newFn func() *tracex.Trace) *tracex.Trace {
	ctx = Get(ctx)

	t, ok := ctx.Value(traceCtxKey).(*tracex.Trace)
	if ok {
		return t
	} else if newFn == nil {
		return nil
	} else {
		return newFn()
	}
}

func WithGinCtx(ctx context.Context, gCtx *gin.Context) context.Context {
	ctx = Get(ctx)
	return context.WithValue(ctx, ginCtxKey, gCtx)
}

func GinCtxFromCtxOrNew(ctx context.Context, newFu func() *gin.Context) *gin.Context {
	if gCtx, ok := ctx.(*gin.Context); ok {
		return gCtx
	}

	gCtx, ok := ctx.Value(ginCtxKey).(*gin.Context)
	if ok {
		return gCtx
	} else if newFu == nil {
		return nil
	} else {
		return gCtx
	}
}

func Get(c context.Context) context.Context {
	if gCtx := GinCtxFromCtxOrNew(c, nil); gCtx != nil {
		c = gCtx.Request.Context()
	}
	return c
}

func Copy(ctx context.Context) context.Context {
	t := TraceFromCtxOrNew(ctx, nil)
	return WithTraceCtx(context.Background(), t)
}

func New(ctx context.Context, req *http.Request) context.Context {
	if _, ok := ctx.(*gin.Context); ok {
		ginC := ctx.(*gin.Context)
		ctx = ginC.Request.Context()
	}
	if v := ctx.Value(traceCtxKey); v == nil {
		ctx = WithTraceCtx(ctx, tracex.New(req))
	}
	return ctx
}
