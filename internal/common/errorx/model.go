package errorx

type CodeEntry struct {
	Code    int
	Message string
}

type ErrMeta struct {
	IsExternalErr   bool       `json:"-"`
	ExternalErrType *CodeEntry `json:"external_err_type,omitempty"`

	ServiceType *CodeEntry `json:"service_type,omitempty"`
	Service     *CodeEntry `json:"service,omitempty"`
	ErrType     *CodeEntry `json:"err_type,omitempty"`
	IsSuccess   *CodeEntry `json:"is_success,omitempty"`
	InnerCode   *CodeEntry `json:"inner_code,omitempty"`
}

type MErr struct {
	ErrMeta
	IsFinalSuccess bool   `json:"is_final_success"`
	FinalCode      int    `json:"final_code"`
	FinalMsg       string `json:"final_msg"`
}

type ErrOptions struct {
	ErrMeta
	Code int
	Err  error
}
