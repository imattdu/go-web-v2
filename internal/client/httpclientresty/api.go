package httpclientresty

import (
	"context"
	"net/http"
)

func Post(ctx context.Context, req *Req) error {
	req.Meta.method = http.MethodPost
	return do(ctx, req)
}

func Get(ctx context.Context, req *Req) error {
	req.Meta.method = http.MethodGet
	return do(ctx, req)
}
