package errorx

import (
	"errors"
	"fmt"
	"strconv"
)

func (e MErr) Error() string {
	return e.FinalMsg
}

func New(query ErrOptions) error {
	if query.Err == nil {
		return nil
	}
	// 1.优先级：参数 > 参数中的err > 默认兜底
	rsp := MErr{
		ErrMeta:  query.ErrMeta,
		FinalMsg: query.Err.Error(),
	}
	if query.Code != 0 {
		rsp.InnerCode = &CodeEntry{
			Code:    query.Code,
			Message: query.Err.Error(),
		}
	}

	// 2.mErr 优先级
	var mErr MErr
	if errors.As(query.Err, &mErr) {
		if rsp.IsExternalErr {
			if rsp.ExternalErrType == nil && (mErr.ServiceType == ServiceTypeBasic || mErr.ServiceType == ServiceTypeService) {
				rsp.ExternalErrType = ExternalErrTypeService
			}
		} else {
			if rsp.ErrType == nil {
				rsp.ErrType = mErr.ErrType
			}
		}
		if rsp.IsSuccess == nil {
			rsp.IsSuccess = mErr.IsSuccess
		}
		if rsp.InnerCode == nil {
			rsp.InnerCode = mErr.InnerCode
		}
	}

	// 3.默认兜底逻辑
	if rsp.IsExternalErr {
		if rsp.ExternalErrType == nil {
			rsp.ExternalErrType = ExternalErrTypeDefault
		}
	} else {
		if rsp.ServiceType == nil {
			rsp.ServiceType = ServiceTypeDefault
		}
		if rsp.ErrType == nil {
			rsp.ErrType = ErrTypeSys
		}
	}
	if rsp.Service == nil {
		rsp.Service = ServiceDefault
	}
	if rsp.IsSuccess == nil {
		rsp.IsSuccess = Failed
	}
	if rsp.InnerCode == nil {
		rsp.InnerCode = ErrDefault
	}
	rsp.FinalCode = rsp.code()
	rsp.IsFinalSuccess = rsp.IsSuccess == Success
	return rsp
}

func (e MErr) code() int {
	arr := make([]int, 0, 1)
	if e.IsExternalErr {
		arr = []int{e.ExternalErrType.Code, e.Service.Code, e.IsSuccess.Code, e.InnerCode.Code}
	} else {
		arr = []int{e.ServiceType.Code, e.Service.Code, e.ErrType.Code, e.IsSuccess.Code, e.InnerCode.Code}
	}
	var codeStr string
	for idx, v := range arr {
		f := "%s%d"
		if idx == 1 {
			f = "%s%03d"
		}
		codeStr = fmt.Sprintf(f, codeStr, v)
	}
	code, _ := strconv.Atoi(codeStr)
	return code
}

func Get(err error, isExternalErr bool) *MErr {
	if err == nil {
		return nil
	}
	var rsp MErr
	if errors.As(err, &rsp) && isExternalErr == rsp.IsExternalErr {
		return &rsp
	}
	return Get(New(ErrOptions{
		ErrMeta: ErrMeta{
			IsExternalErr: isExternalErr,
		},
		Err: err,
	}), isExternalErr)
}
