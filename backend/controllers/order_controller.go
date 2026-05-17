package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	if err := config.DB.Where("id IN ? AND user_id = ?", req.CartIDs, req.UserID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取购物车失败"})
		return
	}

	if len(carts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择有效商品"})
		return
	}

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

	orderNo := "ORD" + time.Now().Format("20060102150405") + time.Now().Format("999999")

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

	for i := range items {
		items[i].OrderID = order.ID
		config.DB.Create(&items[i])
	}

	config.DB.Unscoped().Delete(&carts)

	c.JSON(http.StatusOK, gin.H{
		"message": "订单创建成功",
		"order":   order,
	})
}

// 获取订单列表
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

func PayOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	// 1. 查询订单
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	// 2. 防止重复支付
	if order.Status == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单已支付"})
		return
	}

	// 3. 开启事务（所有操作要么全成功，要么全失败）
	tx := config.DB.Begin()

	// 4. 查询订单商品
	var orderItems []models.OrderItem
	if err := tx.Where("order_id = ?", id).Find(&orderItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单商品失败"})
		return
	}

	// 5. 扣减库存（绝对安全，不会超卖）
	for _, item := range orderItems {
		err := tx.Model(&models.Product{}).
			Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
			UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "库存不足，支付失败"})
			return
		}
	}

	// 6. 更新订单状态
	if err := tx.Model(&order).UpdateColumn("status", 1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "订单状态更新失败"})
		return
	}

	// 7. 提交事务
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "支付成功"})
}

// 获取订单商品
func GetOrderItems(c *gin.Context) {
	orderID := c.Param("id")
	var items []models.OrderItem

	config.DB.Where("order_id = ?", orderID).Find(&items)

	type Item struct {
		Name string `json:"name"`
	}
	var res []Item
	for _, v := range items {
		var p models.Product
		config.DB.First(&p, v.ProductID)
		res = append(res, Item{Name: p.Name})
	}

	c.JSON(http.StatusOK, res)
}

// 删除订单
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	config.DB.Unscoped().Delete(&models.Order{}, id)
	config.DB.Unscoped().Where("order_id = ?", id).Delete(&models.OrderItem{})
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
