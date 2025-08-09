package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *PingHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", h.Ping)
}
