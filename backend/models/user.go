package models

import (
	"time"

	"gorm.io/gorm"
)

// 定义结构体
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"` // 必须有 gorm:"not null"
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `gorm:"default:1" json:"status"`
}
