package httpclient

import (
	"net/http"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/errorx"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

// callStats records timing and raw response info.
type callStats struct {
	attempt   int
	lastError error
	rpcFinal  bool
}

type HttpRequest struct {
	URL          string            `json:"url"`
	QueryParams  map[string]string `json:"query_params"`
	Headers      map[string]string `json:"headers"`
	JSONBody     any               `json:"request_body,omitempty"`
	ResponseBody any               `json:"response_body,omitempty"`

	Timeout time.Duration `json:"timeout"`
	Retries int           `json:"retry_count"`
	RetryIf func(*http.Response, error) bool
	IsError func(*http.Response) error

	Service *errorx.CodeEntry
	Stats   callStats
	method  string
}
