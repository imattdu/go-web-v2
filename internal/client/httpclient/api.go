package httpclient

import (
	"context"
	"net/http"
)

func Get(ctx context.Context, req *Req) (err error) {
	req.Meta.Method = http.MethodGet
	return do(ctx, req)
}

func Post(ctx context.Context, req *Req) (err error) {
	req.Meta.Method = http.MethodPost
	return do(ctx, req)
}
