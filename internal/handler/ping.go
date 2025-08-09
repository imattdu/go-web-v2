package handler

import (
	"github.com/imattdu/go-web-v2/internal/render"
	"github.com/imattdu/go-web-v2/internal/service"
	"github.com/imattdu/go-web-v2/internal/util/logger"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	service.PingService
}

func NewPingHandler() *PingHandler {
	return &PingHandler{
		PingService: service.NewPingService(),
	}
}

func (h *PingHandler) Ping(c *gin.Context) {
	logger.Info(c, logger.TagUndef, map[string]interface{}{
		"msg": "pong",
	})
	render.JSON(c, 200, map[string]interface{}{
		"ping": h.PingService.Ping(c, service.PingParams{
			Name: c.Param("name"),
		}),
	}, nil)
}
