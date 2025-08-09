package router

import (
	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/handler/ping"
	"github.com/imattdu/go-web-v2/internal/handler/user"
)

func RegisterRoutes(r *gin.Engine) {
	group := r.Group("/go-web-v2")
	ping.NewHandler().RegisterRoutes(group)
	user.NewHandler().RegisterRoutes(group)
}
