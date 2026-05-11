package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 添加到购物车
func AddToCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existingCart models.Cart
	config.DB.Where("user_id = ? AND product_id = ?", cart.UserID, cart.ProductID).First(&existingCart)
	if existingCart.ID != 0 {
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
	config.DB.Delete(&models.Cart{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
