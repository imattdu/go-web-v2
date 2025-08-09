package api

import "github.com/imattdu/go-web-v2/internal/model"

type ListResult struct {
	Users []model.User `json:"users"`
}
