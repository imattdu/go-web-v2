package api

import (
	"context"

	"github.com/imattdu/go-web-v2/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, params ListByNameParams, tx *gorm.DB) ([]model.User, error)
	Get(ctx context.Context, tx *gorm.DB) (model.User, error)
}
