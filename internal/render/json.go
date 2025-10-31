package render

import (
	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/common/cctxv2"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/trace"
	"net/http"
)

func JSON(ctx *gin.Context, status int, data interface{}, err error) {
	c := ctx.Request.Context()
	t, ok := cctxv2.GetAs[*trace.Trace](c, cctxv2.TraceKey)
	if !ok {
		t = trace.New(&http.Request{})
	}

	var (
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
