package httptrace

import (
	"bytes"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/trace"
	"github.com/imattdu/go-web-v2/internal/common/util"
	logger2 "github.com/imattdu/go-web-v2/internal/common/util/logger"
	"github.com/imattdu/go-web-v2/internal/render"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	body *bytes.Buffer
	gin.ResponseWriter
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Req() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			t    = trace.New(ctx.Request)
			req  = ctx.Request
			dCtx = cctx.WithGinCtx(req.Context(), ctx)
		)
		dCtx = cctx.WithTraceCtx(ctx, t)
		req = req.WithContext(dCtx)
		ctx.Request = req

		// 获取body
		reqBodyBytes, err := ctx.GetRawData()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			logger2.Warn(ctx, logger2.TagUndef, map[string]interface{}{
				"msg": "GetRawData failed",
				"err": err.Error(),
			})
			return
		}
		// 重置HTTP请求体的偏移量
		ctx.Request.Body = io.NopCloser(bytes.NewReader(reqBodyBytes))
		var reqBody interface{}
		_ = util.Unmarshal(ctx, string(reqBodyBytes), &reqBody)
		logger2.Info(ctx, logger2.TagRequestIn, map[string]interface{}{
			logger2.KRequestBody: reqBody,
		})

		// 捕捉响应
		responseWriter := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = responseWriter
		start := time.Now()
		ctx.Next()
		var (
			latency = time.Since(start).Milliseconds()
			// 获取响应 body
			rspBodyStr = responseWriter.body.String()
			rspBody    render.Response
		)
		if err := util.Unmarshal(ctx, rspBodyStr, &rspBody); err != nil {
			return
		}
		logMap := map[string]interface{}{
			logger2.KCode:         rspBody.Code,
			logger2.KRequestBody:  reqBody,
			logger2.KResponseBody: rspBody,
			logger2.KProcTime:     latency,
		}
		if rspBody.ErrType != 0 && rspBody.ErrType != errorx.ErrTypeBiz.Code {
			logger2.Warn(ctx, logger2.TagRequestOut, logMap)
			return
		}
		logger2.Info(ctx, logger2.TagRequestOut, logMap)
		rspBody.ErrType = -1
	}
}
