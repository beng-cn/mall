package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var req struct {
		CartIDs []uint `json:"cart_ids"`
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

	var carts []models.Cart
	if err := config.DB.Where("id IN ? AND user_id = ?", req.CartIDs, userID.(uint)).Find(&carts).Error; err != nil {
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
		UserID:  userID.(uint),
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

	config.DB.Delete(&carts)
	c.JSON(http.StatusOK, gin.H{"message": "订单创建成功", "order": order})
}

func GetOrderList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录，请先登录"})
		return
	}

	var orders []models.Order
	err := config.DB.Where("user_id = ?", userID.(uint)).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单失败"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func PayOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权支付"})
		return
	}

	if order.Status == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单已支付"})
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统繁忙"})
		return
	}

	var orderItems []models.OrderItem
	if err := tx.Where("order_id = ? AND deleted_at IS NULL", id).Find(&orderItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单商品失败"})
		return
	}

	// 扣减库存
	for _, item := range orderItems {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "商品不存在"})
			return
		}

		if product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "商品库存不足"})
			return
		}

		tx.Model(&product).UpdateColumn("stock", product.Stock-item.Quantity)

		ClearSingleProductCache(fmt.Sprintf("%d", item.ProductID))
	}

	tx.Model(&order).UpdateColumn("status", 1)
	tx.Commit()

	ClearAllProductListCache()

	c.JSON(http.StatusOK, gin.H{"message": "支付成功"})
}

func GetOrderItems(c *gin.Context) {
	orderID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var order models.Order
	if err := config.DB.Where("id = ? AND user_id = ?", orderID, userID.(uint)).First(&order).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权查看"})
		return
	}

	var items []models.OrderItem
	config.DB.Where("order_id = ?", orderID).Find(&items)

	type Item struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}
	var res []Item
	for _, v := range items {
		var p models.Product
		config.DB.First(&p, v.ProductID)
		res = append(res, Item{
			Name:     p.Name,
			Price:    v.Price,
			Quantity: v.Quantity,
		})
	}

	c.JSON(http.StatusOK, res)
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var order models.Order
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID.(uint)).First(&order).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除"})
		return
	}

	config.DB.Delete(&order)
	config.DB.Where("order_id = ?", id).Delete(&models.OrderItem{})
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
