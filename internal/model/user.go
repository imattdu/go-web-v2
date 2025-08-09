package model

import "time"

// User  用户表
type User struct {
	ID           int64     `gorm:"column(id)" json:"id"`                       //  用户id，主键
	Username     string    `gorm:"column(username)" json:"username"`           //  用户名
	Email        string    `gorm:"column(email)" json:"email"`                 //  邮箱，唯一
	PasswordHash string    `gorm:"column(password_hash)" json:"password_hash"` //  密码哈希值
	CreatedAt    time.Time `gorm:"column(created_at)" json:"created_at"`       //  创建时间
	UpdatedAt    time.Time `gorm:"column(updated_at)" json:"updated_at"`       //  更新时间
	Status       int       `gorm:"column(status)" json:"status"`
}

func (User) TableName() string {
	return "user"
}
