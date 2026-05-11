package models

import "gorm.io/gorm"

//定义结构体
type Order struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	OrderNo string  `gorm:"unique;not null" json:"order_no"`
	Total   float64 `json:"total"`
	Status  int     `gorm:"default:0" json:"status"` // 0待支付 1已支付 2已取消
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
