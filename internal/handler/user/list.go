package user

import (
	userConverter "github.com/imattdu/go-web-v2/internal/converter/user"
	userDTO "github.com/imattdu/go-web-v2/internal/dto/user"
	"github.com/imattdu/go-web-v2/internal/render"
	userService "github.com/imattdu/go-web-v2/internal/service/user"
	userAPI "github.com/imattdu/go-web-v2/internal/service/user/api"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService userAPI.Service
}

func NewHandler() *Handler {
	return &Handler{
		UserService: userService.NewService(),
	}
}

func (h *Handler) ListUser(c *gin.Context) {
	var req userDTO.ListRequest
	if err := c.BindJSON(&req); err != nil {
		render.JSON(c, 200, nil, err)
		return
	}

	params, err := userConverter.ListRequestToParams(c, req)
	if err != nil {
		render.JSON(c, 200, nil, err)
		return
	}
	result, err := h.UserService.List(c.Request.Context(), params)
	if err != nil {
		render.JSON(c, 200, nil, err)
		return
	}
	response, err := userConverter.ListParamsToResponse(c.Request.Context(), result)
	render.JSON(c, 200, response, err)
	return
}
