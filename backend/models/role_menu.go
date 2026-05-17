package models

import (
	"time"

	"gorm.io/gorm"
)

// 定义结构体
type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name string `gorm:"not null" json:"name"`
	Desc string `json:"desc"`
}

type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name     string `json:"name"`
	ParentID uint   `gorm:"default:0" json:"parent_id"`
	Path     string `json:"path"`
}

type RoleMenu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	RoleID uint `json:"role_id"`
	MenuID uint `json:"menu_id"`
}
