package user

import (
	"context"
	"time"

	userRepo "github.com/imattdu/go-web-v2/internal/repository/user"
	userRepoAPI "github.com/imattdu/go-web-v2/internal/repository/user/api"
	userAPI "github.com/imattdu/go-web-v2/internal/service/user/api"
)

type service struct {
	userRepo userRepoAPI.Repository
}

func NewService() userAPI.Service {
	return &service{
		userRepo: userRepo.NewRepository(),
	}
}

func (s *service) List(ctx context.Context, params userAPI.ListParams) (userAPI.ListResult, error) {
	time.Sleep(time.Millisecond * 100)
	rsp, err := s.userRepo.List(ctx, userRepoAPI.ListByNameParams{
		Username: params.Username,
	}, nil)
	if err != nil {
		return userAPI.ListResult{}, err
	}
	return userAPI.ListResult{
		Users: rsp,
	}, nil
}
