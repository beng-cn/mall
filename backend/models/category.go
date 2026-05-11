package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	ParentID uint   `gorm:"default:0" json:"parent_id"`
	Status   int    `gorm:"default:1" json:"status"`
}
