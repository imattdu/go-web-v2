package errorx

var (
	ExternalErrTypeDefault = &CodeEntry{
		Code:    1,
		Message: "未知错误",
	}
	ExternalErrTypeSys = &CodeEntry{
		Code:    2,
		Message: "系统错误",
	}
	ExternalErrTypeService = &CodeEntry{
		Code:    4,
		Message: "服务错误",
	}
	ExternalErrTypeBiz = &CodeEntry{
		Code:    5,
		Message: "业务错误",
	}
)

var (
	ServiceTypeDefault = &CodeEntry{
		Code:    1,
		Message: "未知服务类型",
	}
	ServiceTypeBasic = &CodeEntry{
		Code:    4,
		Message: "组件",
	}
	ServiceTypeService = &CodeEntry{
		Code:    5,
		Message: "服务",
	}
)

var (
	ServiceDefault = &CodeEntry{
		Code:    1,
		Message: "未知服务",
	}
	ServiceMysql = &CodeEntry{
		Code:    4,
		Message: "mysql",
	}
	ServiceRedis = &CodeEntry{
		Code:    6,
		Message: "redis",
	}

	ServiceUser = &CodeEntry{
		Code:    19,
		Message: "user",
	}
)

var (
	ErrTypeSys = &CodeEntry{
		Code:    4,
		Message: "系统错误",
	}
	ErrTypeBiz = &CodeEntry{
		Code:    4,
		Message: "业务错误",
	}
)

var (
	Success = &CodeEntry{
		Code:    0,
		Message: "success",
	}
	Failed = &CodeEntry{
		Code:    1,
		Message: "failed",
	}
)

var (
	ErrDefault = &CodeEntry{
		Code:    1000,
		Message: "未知错误",
	}
	ErrNotFound = &CodeEntry{
		Code:    404,
		Message: "not found",
	}
)
