package logger

const (
	logK = "logger"

	KErr    = "err"
	KErrMsg = "err_msg"

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
	TagRedisSuccess = "_com_redis_success"
	TagRedisFailure = "_com_redis_failure"

	KAttempt      = "attempt"
	KRetries      = "retries"
	KIsFinal      = "is_final"
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

	KCallee     = "callee"
	KCalleeFunc = "callee_func"

	KParams = "params"

	KSql = "sql"
	KCmd = "cmd"
)
