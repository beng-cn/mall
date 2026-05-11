package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"time"

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

	// 查询购物车
	var carts []models.Cart
	if err := config.DB.Where("id IN ? AND user_id = ?", req.CartIDs, req.UserID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取购物车失败"})
		return
	}

	if len(carts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择有效商品"})
		return
	}

	// 计算总价 + 组装订单项
	var total float64
	var items []models.OrderItem
	for _, cart := range carts {
		var product models.Product
		if err := config.DB.First(&product, cart.ProductID).Error; err != nil {
			continue
		}
		total += product.Price * float64(cart.Quantity)
		items = append(items, models.OrderItem{
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
			Price:     product.Price,
		})
	}

	// 生成唯一订单号（时间戳 + 纳秒）
	orderNo := "ORD" + time.Now().Format("20060102150405") + time.Now().Format("999999")

	// 创建订单
	order := models.Order{
		UserID:  req.UserID,
		OrderNo: orderNo,
		Total:   total,
		Status:  0,
	}
	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 创建订单明细
	for i := range items {
		items[i].OrderID = order.ID
		config.DB.Create(&items[i])
	}

	// 删除已结算购物车
	config.DB.Delete(&carts)

	c.JSON(http.StatusOK, gin.H{
		"message": "订单创建成功",
		"order":   order,
	})
}

// ==============================
// 新增：获取当前用户的订单列表（给我的订单页面使用）
// ==============================
func GetOrderList(c *gin.Context) {
	userID := c.Query("user_id")

	var orders []models.Order
	err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单失败"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
