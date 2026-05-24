package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProductList(c *gin.Context) {
	keyword := c.Query("keyword")
	categoryID := c.Query("category_id")

	cacheKey := fmt.Sprintf("product:list:%s:%s", keyword, categoryID)

	cacheData, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，直接返回
		var products []models.Product
		json.Unmarshal([]byte(cacheData), &products)
		c.JSON(200, products)
		return
	}

	// 3. 缓存没命中，查数据库
	var products []models.Product
	db := config.DB

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	if err := db.Find(&products).Error; err != nil {
		c.JSON(500, gin.H{"error": "查询失败"})
		return
	}

	// 4. 存入缓存，5分钟过期
	data, _ := json.Marshal(products)
	config.RDB.Set(config.Ctx, cacheKey, data, 5*time.Minute)

	c.JSON(200, products)
}

// 获取商品详情
func GetProductDetail(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "product:item:" + id

	// 读缓存
	cacheData, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var product models.Product
		json.Unmarshal([]byte(cacheData), &product)
		c.JSON(200, product)
		return
	}

	// 查库
	var product models.Product
	config.DB.First(&product, id)

	// 写缓存
	data, _ := json.Marshal(product)
	config.RDB.Set(config.Ctx, cacheKey, data, 10*time.Minute)

	c.JSON(200, product)
}

// 获取所有父分类
func GetParentCategories(c *gin.Context) {
	cacheKey := "category:parent"

	// 1. 先查Redis
	cacheData, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var list []models.Category
		_ = json.Unmarshal([]byte(cacheData), &list)
		c.JSON(http.StatusOK, list)
		return
	}

	// 2. 缓存没命中，执行你原来的数据库查询
	var list []models.Category
	config.DB.Where("parent_id = 0 AND status = 1").Find(&list)

	// 3. 写入缓存，1小时过期（分类很少改动）
	data, _ := json.Marshal(list)
	config.RDB.Set(config.Ctx, cacheKey, data, 1*time.Hour)

	c.JSON(http.StatusOK, list)
}

// 获取子分类
func GetChildCategories(c *gin.Context) {
	parentID := c.Query("parent_id")
	cacheKey := fmt.Sprintf("category:child:%s", parentID)

	// 1. 查缓存
	cacheData, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var list []models.Category
		_ = json.Unmarshal([]byte(cacheData), &list)
		c.JSON(http.StatusOK, list)
		return
	}

	var list []models.Category
	config.DB.Where("parent_id = ? AND status = 1", parentID).Find(&list)

	// 3. 写缓存
	data, _ := json.Marshal(list)
	config.RDB.Set(config.Ctx, cacheKey, data, 1*time.Hour)

	c.JSON(http.StatusOK, list)
}
