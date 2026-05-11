package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 全局跨域中间件
	r.Use(middlewares.Cors())

	// 用户端路由
	user := r.Group("/api/user")
	{
		user.POST("/register", controllers.Register)
		user.POST("/login", controllers.Login)
		user.PUT("/info", controllers.UpdateUserInfo)
	}

	product := r.Group("/api/product")
	{
		product.GET("/list", controllers.GetProductList)
		product.GET("/:id", controllers.GetProductDetail)
	}

	cart := r.Group("/api/cart")
	{
		cart.POST("/add", controllers.AddToCart)
		cart.PUT("/:id", controllers.UpdateCartQuantity)
		cart.DELETE("/:id", controllers.DeleteCartItem)
	}

	order := r.Group("/api/order")
	{
		order.POST("/create", controllers.CreateOrder)
	}

	// 管理端路由
	admin := r.Group("/api/admin")
	{
		admin.POST("/product", controllers.CreateProduct)
	}

	return r
}
