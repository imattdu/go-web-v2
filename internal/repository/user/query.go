package user

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/database/mysql"
	"github.com/imattdu/go-web-v2/internal/model"
	userAPI "github.com/imattdu/go-web-v2/internal/repository/user/api"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository() userAPI.Repository {
	return &repository{
		db: mysql.GlobalDb,
	}
}

func (r *repository) List(ctx context.Context, params userAPI.ListByNameParams, tx *gorm.DB) ([]model.User, error) {
	if tx == nil {
		tx = r.db
	}
	var users []model.User
	err := tx.WithContext(mysql.SetCallStatsClone(ctx, &mysql.CallStats{
		Params: params,
	})).Where("username = ?", params.Username).Find(&users).Error
	return users, err
}

func (r *repository) Get(ctx context.Context, tx *gorm.DB) (model.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user model.User
	err := tx.WithContext(mysql.SetCallStatsClone(ctx, &mysql.CallStats{})).
		Where("username = ?", "haha").
		First(&user).Error
	return user, err
}
