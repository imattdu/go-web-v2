package main

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/client/httpclient"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"github.com/imattdu/go-web-v2/internal/common/config"
	"github.com/imattdu/go-web-v2/internal/common/logger"
	"github.com/imattdu/go-web-v2/internal/common/trace"
	"github.com/imattdu/go-web-v2/internal/database/mysql"
	"github.com/imattdu/go-web-v2/internal/database/redis"
	"github.com/imattdu/go-web-v2/internal/middleware/httptrace"
	"github.com/imattdu/go-web-v2/internal/router"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Init(ctx context.Context) error {
	if err := config.Init("./config/local.toml"); err != nil {
		return err
	}
	logger.Init()
	if err := mysql.Init(ctx); err != nil {
		return err
	}

	httpclient.Init()
	redis.Init()
	return nil
}

func main() {
	ctx := cctx.New(context.Background(), map[string]any{})
	trace.SetTrace(ctx, trace.New(&http.Request{
		URL: &url.URL{},
	}))
	if err := Init(ctx); err != nil {
		return
	}

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	r := gin.New()
	// 绑定中间件 需要放在mapping前面
	r.Use(httptrace.Req())

	router.RegisterRoutes(r)
	if err := r.Run(":8001"); err != nil {
		logger.Info(ctx, logger.TagUndef, map[string]interface{}{
			"err": err.Error(),
			"msg": "server start failed",
		})
	}
}
