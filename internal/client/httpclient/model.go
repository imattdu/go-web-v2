package httpclient

import (
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

// ReqMeta carries request and response metadata.
type ReqMeta struct {
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	URLParams    map[string]string `json:"url_params"`
	Headers      map[string]string
	RequestBody  interface{} `json:"request_body"`
	ResponseBody interface{} `json:"response_body"`

	Timeout    time.Duration `json:"timeout"`
	RetryCount int           `json:"retry_count"`
	RetryIf    func(*http.Response, []byte, error) bool
	OnError    func(*http.Response, []byte) error
}

// callStats records timing and raw response info.
type callStats struct {
	startTime    time.Time
	duration     time.Duration
	rawResponse  *http.Response
	responseText []byte
	errs         []error
	code         int

	retry      int
	isRpcFinal bool
}

// Req wraps the full HTTP request context.
type Req struct {
	Service *errorx.CodeEntry
	Meta    ReqMeta
	Stats   callStats
	client  *gorequest.SuperAgent
}
