package render

import (
	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
)

func JSON(ctx *gin.Context, status int, data interface{}, err error) {
	var (
		trace = cctx.TraceFromCtxOrNew(ctx, nil)
		mErr  = errorx.Get(err, true)
		rsp   = Response{
			TraceId: trace.TraceId.V(),
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
