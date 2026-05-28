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

	// ==================== 公开接口（完全不需要登录） ====================
	user := r.Group("/api/user")
	{
		user.POST("/register", controllers.Register)
		user.POST("/login", controllers.Login)
	}

	product := r.Group("/api/product")
	{
		product.GET("/list", controllers.GetProductList)
		product.GET("/:id", controllers.GetProductDetail)
		product.GET("/category/parents", controllers.GetParentCategories)
		product.GET("/category/children", controllers.GetChildCategories)
	}

	// ✅ 支付宝相关公开接口
	alipay := r.Group("/api/alipay")
	{
		// 异步回调（生产环境用，本地支付宝无法访问）
		alipay.POST("/notify", controllers.AliPayNotify)
		// 同步跳转（本地开发核心，支付成功后自动跳转这里）
		alipay.GET("/success", controllers.AliPayReturn)
	}

	// ==================== 需要普通用户登录的接口 ====================
	auth := r.Group("/api/auth")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.PUT("/user/info", controllers.UpdateUserInfo)
		cart := auth.Group("/cart")
		{
			cart.GET("/list", controllers.GetCartList)
			cart.POST("/add", controllers.AddToCart)
			cart.PUT("/:id", controllers.UpdateCartQuantity)
			cart.DELETE("/:id", controllers.DeleteCartItem)
		}
		order := auth.Group("/order")
		{
			order.POST("/create", controllers.CreateOrder)
			order.GET("/list", controllers.GetOrderList)
			// 支付宝统一下单（获取支付链接）
			order.POST("/alipay", controllers.AliPayUnifiedOrder)
			// 保留原模拟支付接口，方便不登录支付宝时测试
			order.POST("/pay/:id", controllers.PayOrder)
			order.GET("/items/:id", controllers.GetOrderItems)
			order.DELETE("/delete/:id", controllers.DeleteOrder)
		}
	}

	// ==================== 需要管理员权限的接口 ====================
	admin := r.Group("/api/admin")
	admin.Use(middlewares.AdminAuthMiddleware())
	{
		// 商品管理
		admin.POST("/product", controllers.CreateProduct)
		admin.PUT("/product/:id", controllers.UpdateProduct)
		admin.DELETE("/product/:id", controllers.DeleteProduct)
		// 分类管理
		admin.POST("/category/add", controllers.CreateCategory)
		admin.PUT("/category/:id", controllers.UpdateCategory)
		admin.DELETE("/category/:id", controllers.DeleteCategory)
		// 图片上传
		admin.POST("/upload", controllers.UploadImage)
		// 用户管理
		admin.GET("/user/list", controllers.ListUsers)
		admin.PUT("/user/:id/status", controllers.UpdateUserStatus)
		admin.DELETE("/user/:id", controllers.DeleteUser)
	}

	return r
}
