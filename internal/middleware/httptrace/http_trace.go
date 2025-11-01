package httptrace

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/logger"
	"github.com/imattdu/go-web-v2/internal/common/trace"
	"github.com/imattdu/go-web-v2/internal/common/util"
	"github.com/imattdu/go-web-v2/internal/render"

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
			t   = trace.New(ctx.Request)
			req = ctx.Request
			c   = cctx.New(context.Background(), map[string]any{
				cctx.TraceKey: t,
			})
		)
		req = req.WithContext(c)
		ctx.Request = req

		// 获取body
		reqBodyBytes, err := ctx.GetRawData()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			logger.Warn(c, logger.TagUndef, map[string]interface{}{
				"msg": "GetRawData failed",
				"err": err.Error(),
			})
			return
		}
		// 重置HTTP请求体的偏移量
		ctx.Request.Body = io.NopCloser(bytes.NewReader(reqBodyBytes))
		var reqBody interface{}
		_ = util.Unmarshal(ctx, string(reqBodyBytes), &reqBody)
		logger.Info(c, logger.TagRequestIn, map[string]interface{}{
			logger.KRequestBody: reqBody,
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
			logger.KRequestBody:  reqBody,
			logger.KResponseBody: rspBody,
			logger.KProcTime:     latency,
		}
		if rspBody.ErrType != 0 && rspBody.ErrType != errorx.ErrTypeBiz.Code {
			logger.Warn(c, logger.TagRequestOut, logMap)
			return
		}
		logger.Info(c, logger.TagRequestOut, logMap)
		rspBody.ErrType = -1
	}
}
