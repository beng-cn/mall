package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
)

// 支付宝统一下单接口（获取支付链接）
func AliPayUnifiedOrder(c *gin.Context) {
	var req struct {
		OrderID uint `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 查询订单
	var order models.Order
	if err := config.DB.Where("id = ? AND user_id = ?", req.OrderID, userID.(uint)).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
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

	// 创建支付宝支付请求
	p := alipay.TradePagePay{}
	p.Subject = fmt.Sprintf("商城订单-%s", order.OrderNo)
	p.OutTradeNo = order.OrderNo                     // 商户订单号，必须唯一
	p.TotalAmount = fmt.Sprintf("%.2f", order.Total) // 支付金额，单位：元
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	p.NotifyURL = config.AliPayNotifyURL
	p.ReturnURL = config.AliPayReturnURL

	// 生成支付链接
	url, err := config.AliPayClient.TradePagePay(p)
	if err != nil {
		log.Println("❌ 生成支付链接失败：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成支付链接失败"})
		return
	}

	log.Println("✅ 支付宝支付链接生成成功：", url.String())
	c.JSON(http.StatusOK, gin.H{
		"message": "生成支付链接成功",
		"pay_url": url.String(),
	})
}

// 支付宝异步回调接口（生产环境用，本地无法访问）
func AliPayNotify(c *gin.Context) {
	log.Println("收到支付宝异步回调请求")
	// 解析回调参数并验证签名
	noti, err := config.AliPayClient.GetTradeNotification(c.Request)
	if err != nil {
		log.Println("❌ 解析异步回调失败：", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	log.Println("✅ 异步回调参数解析成功，订单号：", noti.OutTradeNo, " 交易状态：", noti.TradeStatus)

	// 只处理支付成功的回调
	if noti.TradeStatus != "TRADE_SUCCESS" && noti.TradeStatus != "TRADE_FINISHED" {
		log.Println("⚠️  异步回调非成功状态，忽略：", noti.TradeStatus)
		c.String(http.StatusOK, "success")
		return
	}

	// 根据商户订单号查询订单
	var order models.Order
	if err := config.DB.Where("order_no = ?", noti.OutTradeNo).First(&order).Error; err != nil {
		log.Println("❌ 异步回调未找到订单：", noti.OutTradeNo)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 订单已支付，直接返回成功（幂等性）
	if order.Status == 1 {
		log.Println("✅ 异步回调订单已支付，无需重复处理：", order.OrderNo)
		c.String(http.StatusOK, "success")
		return
	}

	// 调用公共支付处理逻辑
	log.Println("开始处理异步回调订单支付：", order.OrderNo)
	if err := processOrderPayment(order.ID); err != nil {
		log.Println("❌ 异步回调订单处理失败：", order.OrderNo, " 错误：", err)
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	log.Println("✅ 异步回调订单处理成功：", order.OrderNo)
	// 必须返回"success"，否则支付宝会持续回调
	c.String(http.StatusOK, "success")
}

// 支付宝同步跳转接口（本地开发核心，支付成功后会自动跳转到这里）
func AliPayReturn(c *gin.Context) {
	log.Println("\n==================== 收到支付宝同步回调 ====================")
	log.Println("回调URL参数：", c.Request.URL.Query().Encode())

	// 1. 验证签名（关键步骤，失败则后续逻辑不执行）
	log.Println("开始验证签名...")
	err := config.AliPayClient.VerifySign(c.Request.Context(), c.Request.URL.Query())
	if err != nil {
		log.Println("❌ 签名验证失败：", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "支付验证失败：" + err.Error()})
		return
	}
	log.Println("✅ 签名验证通过")

	// 2. 获取订单号
	orderNo := c.Query("out_trade_no")
	log.Println("订单号：", orderNo)

	// 🔥 核心修复：支付宝同步回调 不返回 trade_status，本地开发直接跳过判断
	// 生产环境靠异步回调更新状态，同步回调仅用于本地刷新

	// 3. 查询订单是否存在
	log.Println("查询订单：", orderNo)
	var order models.Order
	if err := config.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		log.Println("❌ 未找到订单：", orderNo, " 错误：", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在：" + orderNo})
		return
	}
	log.Println("✅ 订单查询成功，当前状态：", order.Status)

	// 4. 幂等性判断：订单已支付，直接返回成功
	if order.Status == 1 {
		log.Println("✅ 订单已支付，无需重复处理：", orderNo)
		c.JSON(http.StatusOK, gin.H{
			"message":  "订单已支付",
			"order_id": order.ID,
			"order_no": order.OrderNo,
		})
		return
	}

	// 5. 调用公共支付处理逻辑，更新订单状态+扣减库存（本地开发核心）
	log.Println("开始处理订单支付：", orderNo)
	if err := processOrderPayment(order.ID); err != nil {
		log.Println("❌ 订单支付处理失败：", orderNo, " 错误：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "订单处理失败：" + err.Error()})
		return
	}

	log.Println("✅ 订单支付处理成功：", orderNo)
	c.JSON(http.StatusOK, gin.H{
		"message":  "支付成功",
		"order_id": order.ID,
		"order_no": order.OrderNo,
	})
}

// 公共支付处理逻辑（所有支付方式统一调用）
func processOrderPayment(orderID uint) error {
	log.Println("processOrderPayment 开始处理订单ID：", orderID)

	// 开启事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		log.Println("❌ 开启事务失败：", tx.Error)
		return tx.Error
	}
	log.Println("✅ 事务开启成功")

	// 兜底：发生panic时自动回滚事务
	defer func() {
		if r := recover(); r != nil {
			log.Println("❌ 事务发生panic，回滚事务：", r)
			tx.Rollback()
		}
	}()

	// 1. 查询订单
	var order models.Order
	if err := tx.First(&order, orderID).Error; err != nil {
		log.Println("❌ 查询订单失败：", orderID, " 错误：", err)
		tx.Rollback()
		return err
	}
	log.Println("✅ 订单查询成功，当前状态：", order.Status)

	// 2. 查询订单项
	var orderItems []models.OrderItem
	if err := tx.Where("order_id = ? AND deleted_at IS NULL", orderID).Find(&orderItems).Error; err != nil {
		log.Println("❌ 查询订单项失败：", orderID, " 错误：", err)
		tx.Rollback()
		return err
	}
	log.Println("✅ 订单项查询成功，共", len(orderItems), "个商品")

	// 3. 乐观锁扣减库存（最多重试3次）
	for _, item := range orderItems {
		log.Println("开始扣减商品ID：", item.ProductID, " 数量：", item.Quantity)
		for retry := 0; retry < 3; retry++ {
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				log.Println("❌ 查询商品失败：", item.ProductID, " 错误：", err)
				tx.Rollback()
				return err
			}

			if product.Stock < item.Quantity {
				log.Println("❌ 商品库存不足：", product.Name, " 库存：", product.Stock, " 需要：", item.Quantity)
				tx.Rollback()
				return fmt.Errorf("商品《%s》库存不足，剩余%d件", product.Name, product.Stock)
			}

			// 乐观锁更新库存
			result := tx.Model(&product).
				Where("version = ?", product.Version).
				Updates(map[string]interface{}{
					"stock":   product.Stock - item.Quantity,
					"version": product.Version + 1,
				})

			if result.Error != nil {
				log.Println("❌ 更新库存失败：", product.ID, " 错误：", result.Error)
				tx.Rollback()
				return result.Error
			}

			if result.RowsAffected > 0 {
				log.Println("✅ 商品库存扣减成功：", product.ID)
				// 缓存清理：如果你的项目没有Redis，可以注释掉这行
				ClearSingleProductCache(fmt.Sprintf("%d", item.ProductID))
				break
			}

			log.Println("⚠️  乐观锁更新失败，重试：", retry+1)
			if retry == 2 {
				log.Println("❌ 乐观锁重试次数用完，扣减库存失败")
				tx.Rollback()
				return fmt.Errorf("系统繁忙，请稍后重试")
			}
		}
	}

	// 4. 更新订单状态为已支付
	log.Println("开始更新订单状态为已支付：", orderID)
	if err := tx.Model(&order).UpdateColumn("status", 1).Error; err != nil {
		log.Println("❌ 更新订单状态失败：", orderID, " 错误：", err)
		tx.Rollback()
		return err
	}
	log.Println("✅ 订单状态更新为已支付")

	// 5. 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Println("❌ 事务提交失败：", orderID, " 错误：", err)
		return err
	}
	log.Println("✅ 事务提交成功，订单处理完成")

	// 缓存清理：如果你的项目没有Redis，可以注释掉这行
	ClearAllProductListCache()

	return nil
}
