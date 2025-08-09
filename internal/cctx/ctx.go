package cctx

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/trace"
)

type ctxKey string

var (
	traceCtxKey ctxKey = "traceCtxKey"
	ginCtxKey   ctxKey = "ginCtxKey"
	mysqlCtxKey ctxKey = "mysqlCtxKey"
)

func WithTraceCtx(ctx context.Context, trace *trace.Trace) context.Context {
	ctx = Get(ctx)
	return context.WithValue(ctx, traceCtxKey, trace)
}

func TraceFromCtxOrNew(ctx context.Context, newFn func() *trace.Trace) *trace.Trace {
	ctx = Get(ctx)

	t, ok := ctx.Value(traceCtxKey).(*trace.Trace)
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

type Mysql struct {
	Query interface{} `json:"query"`
	Start time.Time   `json:"start"`
}

func WithMysqlCtx(ctx context.Context, d Mysql) context.Context {
	ctx = Get(ctx)
	return context.WithValue(ctx, mysqlCtxKey, d)
}

func MysqlFromCtx(ctx context.Context) Mysql {
	if v, ok := ctx.Value(mysqlCtxKey).(Mysql); ok {
		return v
	}
	return Mysql{}
}

func Get(c context.Context) context.Context {
	if gCtx := GinCtxFromCtxOrNew(c, nil); gCtx != nil {
		c = gCtx.Request.Context()
	}
	return c
}

func New(ctx context.Context, req *http.Request) context.Context {
	if _, ok := ctx.(*gin.Context); ok {
		ginC := ctx.(*gin.Context)
		ctx = ginC.Request.Context()
	}
	if v := ctx.Value(traceCtxKey); v == nil {
		ctx = WithTraceCtx(ctx, trace.New(req))
	}
	return ctx
}
