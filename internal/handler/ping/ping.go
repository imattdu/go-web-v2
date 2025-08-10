package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/database/redis"
	"github.com/imattdu/go-web-v2/internal/render"
	"github.com/imattdu/go-web-v2/internal/service/ping"
	"time"
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
		"ping": h.PingService.Ping(c.Request.Context()),
	}, nil)
}

type CmdRequest struct {
	Method string        `json:"method"`
	KV     redis.KVEntry `json:"kv"`
}

func (h *Handler) Cmd(c *gin.Context) {
	var req CmdRequest
	if err := c.BindJSON(&req); err != nil {
		render.JSON(c, 200, nil, err)
		return
	}
	req.KV.BaseTTL *= time.Second
	var (
		rsp interface{}
		err error
	)
	switch req.Method {
	case "set":
		err = redis.Set(c.Request.Context(), &req.KV)
	case "get":
		rsp, err = redis.Get(c.Request.Context(), &req.KV)
	case "ttl":
		rsp, err = redis.TTL(c.Request.Context(), &req.KV)
		rspDuration, ok := rsp.(time.Duration)
		if ok {
			rsp = rspDuration / time.Second
		}
	}
	render.JSON(c, 200, rsp, err)
}
