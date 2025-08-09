package ping

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/ping", h.Ping)
}
