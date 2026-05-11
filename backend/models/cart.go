package models

import "gorm.io/gorm"

//定义结构体
type Cart struct {
	gorm.Model
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `gorm:"not null" json:"quantity"`
}
