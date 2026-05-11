package main

import (
	"backend/config"
	"backend/models"
	"backend/routes"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 自动迁移表结构（可选，也可以用SQL脚本建表）
	config.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
		&models.Role{},
		&models.Menu{},
		&models.RoleMenu{},
	)

	// 启动路由
	r := routes.SetupRouter()
	r.Run(":8080")
}
