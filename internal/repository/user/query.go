package user

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/database/mysql"

	"github.com/imattdu/go-web-v2/internal/cctx"
	"github.com/imattdu/go-web-v2/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, params ListByNameParams, tx *gorm.DB) ([]model.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{
		db: mysql.GlobalDb,
	}
}

func (r *repository) List(ctx context.Context, params ListByNameParams, tx *gorm.DB) ([]model.User, error) {
	if tx == nil {
		tx = r.db
	}
	var users []model.User
	err := tx.WithContext(cctx.WithMysqlCtx(ctx, cctx.Mysql{
		Query: params,
	})).Where("username = ?", params.Username).Find(&users).Error
	return users, err
}
