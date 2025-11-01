package user

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {

	group = group.Group("/user")
	group.POST("/list", h.ListUser)
}
