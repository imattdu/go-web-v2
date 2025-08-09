package user

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/model"
	"github.com/imattdu/go-web-v2/internal/repository/user"
)

type Service interface {
	List(ctx context.Context, params ListParams) ([]model.User, error)
}

func NewService() Service {
	return &service{
		UserRepository: user.NewRepository(),
	}
}

type service struct {
	UserRepository user.Repository
}

func (s *service) List(ctx context.Context, params ListParams) ([]model.User, error) {
	return s.UserRepository.List(ctx, user.ListByNameParams{
		Username: params.Username,
	}, nil)
}
