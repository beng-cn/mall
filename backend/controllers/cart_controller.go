package controllers

import (
	"backend/config"
	"backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取购物车列表
func GetCartList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录，请先登录"})
		return
	}

	var carts []models.Cart
	err := config.DB.
		Where("user_id = ? AND deleted_at IS NULL", userID.(uint)).
		Find(&carts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取购物车失败"})
		return
	}

	c.JSON(http.StatusOK, carts)
}

// 添加到购物车
func AddToCart(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录，请先登录"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品不存在"})
		return
	}

	if product.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品库存不足，无法加入购物车"})
		return
	}

	// 查找是否已存在
	var existingCart models.Cart
	err := config.DB.Where(
		"user_id = ? AND product_id = ? AND deleted_at IS NULL",
		userID.(uint), req.ProductID,
	).First(&existingCart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCart := models.Cart{
				UserID:    userID.(uint),
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}
			if err := config.DB.Create(&newCart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "添加购物车失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询购物车失败"})
			return
		}
	} else {
		existingCart.Quantity += req.Quantity
		if err := config.DB.Save(&existingCart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新购物车失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// 修改购物车数量
func UpdateCartQuantity(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 增加权限校验：只能修改自己的购物车
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录，请先登录"})
		return
	}

	// 先查询购物车记录，确认属于当前用户
	var cart models.Cart
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID.(uint)).First(&cart).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改他人的购物车"})
		return
	}

	config.DB.Model(&cart).Update("quantity", req.Quantity)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 删除购物车商品
func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")

	// 增加权限校验：只能删除自己的购物车
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录，请先登录"})
		return
	}

	// 先查询购物车记录，确认属于当前用户
	var cart models.Cart
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID.(uint)).First(&cart).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除他人的购物车"})
		return
	}

	if err := config.DB.Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
