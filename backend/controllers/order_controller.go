package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建订单
func CreateOrder(c *gin.Context) {
	var req struct {
		UserID  uint   `json:"user_id"`
		CartIDs []uint `json:"cart_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var carts []models.Cart
	config.DB.Where("id IN ? AND user_id = ?", req.CartIDs, req.UserID).Find(&carts)
	var total float64
	var items []models.OrderItem
	for _, cart := range carts {
		var product models.Product
		config.DB.First(&product, cart.ProductID)
		total += product.Price * float64(cart.Quantity)
		items = append(items, models.OrderItem{
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
			Price:     product.Price,
		})
	}
	order := models.Order{
		UserID:  req.UserID,
		OrderNo: "ORD" + "123456", // 实际项目需生成唯一订单号
		Total:   total,
		Status:  0,
	}
	config.DB.Create(&order)
	for i := range items {
		items[i].OrderID = order.ID
		config.DB.Create(&items[i])
	}
	config.DB.Delete(&carts)
	c.JSON(http.StatusOK, gin.H{"message": "订单创建成功", "order": order})
}
