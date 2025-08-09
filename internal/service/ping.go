package service

import "github.com/gin-gonic/gin"

type PingService interface {
	Ping(ctx *gin.Context, params PingParams) string
}

type pingService struct{}

func NewPingService() PingService {
	return &pingService{}
}

func (s *pingService) Ping(ctx *gin.Context, params PingParams) string {
	return params.Name
}
