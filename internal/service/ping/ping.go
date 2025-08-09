package ping

import (
	"github.com/gin-gonic/gin"
)

type Service interface {
	Ping(ctx *gin.Context) string
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Ping(ctx *gin.Context) string {
	return "pong"
}
