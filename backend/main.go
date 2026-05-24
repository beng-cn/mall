package main

import (
	"backend/config"
	"backend/models"
	"backend/routes"
)

func main() {
	// 初始化数据库
	config.InitDB()
	config.InitRedis()

	// 自动迁移表结构
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

	r.Static("/uploads", "./uploads")

	r.Run(":8080")
}
