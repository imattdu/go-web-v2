package httpclientresty

import (
	"net/http"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/errorx"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

type ReqMeta struct {
	method       string
	URL          string            `json:"url"`
	URLParams    map[string]string `json:"url_params"`
	Headers      map[string]string
	RequestBody  interface{} `json:"request_body"`
	ResponseBody interface{} `json:"response_body"`

	Timeout    time.Duration `json:"timeout"`
	RetryCount int           `json:"retry_count"`
	RetryIf    func(*http.Response, error) bool
	IsError    func(*http.Response) error
}

// callStats records timing and raw response info.
type callStats struct {
	retry      int
	isRpcFinal bool
}

// Req wraps the full HTTP request context.
type Req struct {
	Service *errorx.CodeEntry
	Meta    ReqMeta
	Stats   callStats
}
