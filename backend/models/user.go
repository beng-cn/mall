package models

import "gorm.io/gorm"

//定义结构体
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"` // 必须有 gorm:"not null"
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `gorm:"default:1" json:"status"`
}
