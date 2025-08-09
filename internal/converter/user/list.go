package user

import (
	"github.com/gin-gonic/gin"
	"github.com/imattdu/go-web-v2/internal/dto/user"
	userService "github.com/imattdu/go-web-v2/internal/service/user"
)

func ListRequestToParams(ctx *gin.Context, request user.ListRequest) userService.ListParams {
	return userService.ListParams{
		Username: request.Username,
	}
}
