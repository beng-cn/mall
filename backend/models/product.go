package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	CategoryID uint    `json:"category_id"`
	Name       string  `gorm:"not null" json:"name"`
	Price      float64 `gorm:"not null" json:"price"`
	Stock      int     `gorm:"not null" json:"stock"`
	Image      string  `json:"image"`
	Status     int     `gorm:"default:1" json:"status"` // 1上架 0下架
}
