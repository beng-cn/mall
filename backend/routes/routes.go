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
		product.GET("/category/parents", controllers.GetParentCategories)
		product.GET("/category/children", controllers.GetChildCategories)
	}

	cart := r.Group("/api/cart")
	{
		cart.GET("/list", controllers.GetCartList)
		cart.POST("/add", controllers.AddToCart)
		cart.PUT("/:id", controllers.UpdateCartQuantity)
		cart.DELETE("/:id", controllers.DeleteCartItem)
	}

	order := r.Group("/api/order")
	{
		order.POST("/create", controllers.CreateOrder)
		order.GET("/list", controllers.GetOrderList)
		order.POST("/pay/:id", controllers.PayOrder)
		order.GET("/items/:id", controllers.GetOrderItems)   // 商品列表
		order.DELETE("/delete/:id", controllers.DeleteOrder) // 删除订单
	}

	// 管理端路由
	admin := r.Group("/api/admin")
	{
		admin.POST("/product", controllers.CreateProduct)
	}

	return r
}
