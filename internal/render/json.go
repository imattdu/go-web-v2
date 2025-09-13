package render

import (
	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/trace"
)

func JSON(ctx *gin.Context, status int, data interface{}, err error) {
	var (
		c    = ctx.Request.Context()
		t    = trace.GetTrace(c)
		mErr = errorx.Get(err, true)
		rsp  = Response{
			TraceId: t.TraceId.V(),
			Data:    data,
			Status:  true,
		}
	)
	if mErr != nil {
		rsp.ErrType = mErr.ExternalErrType.Code
		rsp.Code = mErr.FinalCode
		rsp.Msg = mErr.FinalMsg
		rsp.Status = mErr.IsFinalSuccess
	}
	ctx.JSON(status, rsp)
}
