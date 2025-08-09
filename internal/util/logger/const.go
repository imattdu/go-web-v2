package logger

const (
	logK = "logger"

	KErr = "err"

	KFrom = "from"
	KURI  = "uri"
	KURL  = "url"

	TagUndef        = "_undef"
	TagRequestIn    = "_com_request_in"
	TagRequestOut   = "_com_request_out"
	TagMysqlSuccess = "_com_mysql_success"
	TagMysqlFailure = "_com_mysql_failure"
	TagHttpSuccess  = "_com_http_success"
	TagHttpFailure  = "_com_http_failure"

	KRetry        = "retry"
	KRetryCount   = "retry_count"
	KService      = "service"
	KHeaders      = "headers"
	KRequestBody  = "request_body"
	KResponseBody = "response_body"
	KResponseText = "response_text"
	KIsRPCFinal   = "is_rpc_final"
	KProcTime     = "proc_time"
	KCode         = "code"
)
