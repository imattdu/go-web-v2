package httpclient

import (
	"context"
	"net/http"
)

func Post(ctx context.Context, req *HttpRequest) error {
	req.method = http.MethodPost
	return do(ctx, req)
}

func Get(ctx context.Context, req *HttpRequest) error {
	req.method = http.MethodGet
	return do(ctx, req)
}
