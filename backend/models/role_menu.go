package models

import "gorm.io/gorm"

//定义结构体
type Role struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
	Desc string `json:"desc"`
}

type Menu struct {
	gorm.Model
	Name     string `json:"name"`
	ParentID uint   `gorm:"default:0" json:"parent_id"`
	Path     string `json:"path"`
}

type RoleMenu struct {
	RoleID uint `json:"role_id"`
	MenuID uint `json:"menu_id"`
}
