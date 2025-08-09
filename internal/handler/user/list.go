package user

import (
	"github.com/gin-gonic/gin"
	userConverter "github.com/imattdu/go-web-v2/internal/converter/user"
	userDTO "github.com/imattdu/go-web-v2/internal/dto/user"
	"github.com/imattdu/go-web-v2/internal/render"
	"github.com/imattdu/go-web-v2/internal/service/user"
)

type Handler struct {
	UserService user.Service
}

func NewHandler() *Handler {
	return &Handler{
		UserService: user.NewService(),
	}
}

func (h *Handler) ListUser(c *gin.Context) {
	var req userDTO.ListRequest
	if err := c.BindJSON(&req); err != nil {
		render.JSON(c, 200, nil, err)
		return
	}

	rsp, err := h.UserService.List(c.Request.Context(), userConverter.ListRequestToParams(c, req))
	render.JSON(c, 200, rsp, err)
}
