package httptrace

import (
	"bytes"
	"github.com/imattdu/go-web-v2/internal/cctx"
	"github.com/imattdu/go-web-v2/internal/errorx"
	"github.com/imattdu/go-web-v2/internal/render"
	"github.com/imattdu/go-web-v2/internal/trace"
	"github.com/imattdu/go-web-v2/internal/util"
	"github.com/imattdu/go-web-v2/internal/util/logger"
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
			logger.Warn(ctx, logger.TagUndef, map[string]interface{}{
				"msg": "GetRawData failed",
				"err": err.Error(),
			})
			return
		}
		// 重置HTTP请求体的偏移量
		ctx.Request.Body = io.NopCloser(bytes.NewReader(reqBodyBytes))
		logger.Info(ctx, logger.TagRequestIn, map[string]interface{}{
			logger.KRequestBody: string(reqBodyBytes),
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
			logger.KCode:         rspBody.Code,
			logger.KRequestBody:  string(reqBodyBytes),
			logger.KResponseBody: rspBody,
			logger.KProcTime:     latency,
		}
		if rspBody.ErrType != 0 && rspBody.ErrType != errorx.ErrTypeBiz.Code {
			logger.Warn(ctx, logger.TagRequestOut, logMap)
			return
		}
		logger.Info(ctx, logger.TagRequestOut, logMap)
		rspBody.ErrType = -1
	}
}
