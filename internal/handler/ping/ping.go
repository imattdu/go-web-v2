package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/render"
	"github.com/imattdu/go-web-v2/internal/service/ping"
)

type Handler struct {
	PingService ping.Service
}

func NewHandler() *Handler {
	return &Handler{
		PingService: ping.NewService(),
	}
}

func (h *Handler) Ping(c *gin.Context) {
	render.JSON(c, 200, map[string]interface{}{
		"ping": h.PingService.Ping(c),
	}, nil)
}
