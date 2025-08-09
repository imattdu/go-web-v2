package router

import (
	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/handler"
)

func RegisterRoutes(r *gin.Engine) {
	handler.NewPingHandler().RegisterRoutes(r)
}
