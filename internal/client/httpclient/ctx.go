package httpclient

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
)

var reqCtxKey = "reqCtxKey"

func SetHttpRequest(ctx context.Context, req *HttpRequest) {
	cctx.Set(ctx, reqCtxKey, req)
}

func GetHttpRequest(ctx context.Context) *HttpRequest {
	v, _ := cctx.Get(ctx, reqCtxKey)
	if val, ok := v.(*HttpRequest); ok {
		return val
	}
	return &HttpRequest{}
}
