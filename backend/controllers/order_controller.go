package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 订单创建接口
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

	// 开启数据库事务，保证订单创建全流程原子性
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统繁忙，请稍后再试"})
		return
	}
	// 兜底：发生panic时自动回滚事务，防止数据不一致
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 事务内查询用户选择的购物车商品
	var carts []models.Cart
	if err := tx.Where("id IN ? AND user_id = ?", req.CartIDs, userID.(uint)).Find(&carts).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取购物车失败"})
		return
	}
	if len(carts) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择有效商品"})
		return
	}

	// 2. 计算订单总金额并生成订单项（下单前预校验库存）
	var total float64
	var items []models.OrderItem
	for _, cart := range carts {
		var product models.Product
		if err := tx.First(&product, cart.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("商品ID %d 不存在或已下架", cart.ProductID)})
			return
		}
		// 预校验库存，防止用户下单时库存已不足
		if product.Stock < cart.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("商品《%s》库存不足，剩余%d件", product.Name, product.Stock)})
			return
		}
		total += product.Price * float64(cart.Quantity)
		items = append(items, models.OrderItem{
			ProductID: cart.ProductID,
			Quantity:  cart.Quantity,
			Price:     product.Price,
		})
	}

	// ✅ 修复订单号生成逻辑：使用纳秒时间戳避免高并发下重复
	orderNo := fmt.Sprintf("ORD%s%d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)
	order := models.Order{
		UserID:  userID.(uint),
		OrderNo: orderNo,
		Total:   total,
		Status:  0, // 0=待支付 1=已支付 2=已取消
	}

	// 3. 事务内创建订单
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 4. 批量创建订单项（性能优化：替代循环单个创建）
	for i := range items {
		items[i].OrderID = order.ID
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单项失败"})
		return
	}

	// 5. 事务内删除已下单的购物车商品
	if err := tx.Delete(&carts).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清空购物车失败"})
		return
	}

	// 6. 所有操作成功，提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "订单提交失败，请稍后再试"})
		return
	}

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

// 原模拟支付接口（保留用于测试）
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

	if order.Status == 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单已取消"})
		return
	}

	// 调用公共支付处理逻辑
	if err := processOrderPayment(order.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
