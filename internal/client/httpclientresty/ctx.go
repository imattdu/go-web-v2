package httpclientresty

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
)

type ctxKey string

var reqCtxKey ctxKey = "reqCtxKey"

func WithReqCtx(ctx context.Context, req *Req) context.Context {
	ctx = cctx.Get(ctx)
	return context.WithValue(ctx, reqCtxKey, req)
}

func ReqFromCtx(ctx context.Context) *Req {
	if v, ok := ctx.Value(reqCtxKey).(*Req); ok {
		return v
	}
	return &Req{}
}
