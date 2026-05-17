package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCartList(c *gin.Context) {
	userID := c.Query("user_id") // 从参数获取用户ID

	var carts []models.Cart
	// 查询未删除、属于当前用户的购物车数据
	err := config.DB.
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&carts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	c.JSON(http.StatusOK, carts)
}

// 添加到购物车（增加：库存为0时禁止添加）
func AddToCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ==============================
	// 🔥 新增：先查询商品库存
	// ==============================
	var product models.Product
	if err := config.DB.First(&product, cart.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品不存在"})
		return
	}

	// 🔥 库存为0 → 直接拦截，不让加入购物车
	if product.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品库存不足，无法加入购物车"})
		return
	}

	// 查找是否已存在
	var existingCart models.Cart
	err := config.DB.Where(
		"user_id = ? AND product_id = ? AND deleted_at IS NULL",
		cart.UserID, cart.ProductID,
	).First(&existingCart).Error

	if err == nil {
		existingCart.Quantity += cart.Quantity
		config.DB.Save(&existingCart)
	} else {
		config.DB.Create(&cart)
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// 修改购物车数量
func UpdateCartQuantity(c *gin.Context) {
	id := c.Param("id")
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Model(&models.Cart{}).Where("id = ?", id).Update("quantity", cart.Quantity)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 删除购物车商品
func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Unscoped().Delete(&models.Cart{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
