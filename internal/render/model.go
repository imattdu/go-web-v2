package render

type Response struct {
	TraceId string      `json:"trace_id"`
	ErrType int         `json:"err_type"`
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
