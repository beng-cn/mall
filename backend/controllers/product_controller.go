package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 获取商品列表（支持关键词 + 多分类同时筛选）
func GetProductList(c *gin.Context) {
	var products []models.Product
	keyword := c.Query("keyword")
	categoryIDStr := c.Query("category_id")

	// 🔥 调试日志：打印接收到的参数
	println("=====================================")
	println("接收到的请求参数：")
	println("keyword:", keyword)
	println("category_id:", categoryIDStr)

	db := config.DB

	// 1. 关键词模糊搜索
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
		println("添加关键词筛选：name LIKE %" + keyword + "%")
	}

	// 2. 多分类同时筛选
	if categoryIDStr != "" {
		categoryIDList := strings.Split(categoryIDStr, ",")
		var validCategoryIDs []uint

		for _, idStr := range categoryIDList {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}

			id, err := strconv.ParseUint(idStr, 10, 32)
			if err == nil {
				validCategoryIDs = append(validCategoryIDs, uint(id))
				println("有效分类ID:", id)
			} else {
				println("无效分类ID:", idStr, "错误:", err.Error())
			}
		}

		if len(validCategoryIDs) > 0 {
			// 🔥 开启SQL日志，查看实际执行的查询语句
			db = db.Debug().Where("category_id IN (?)", validCategoryIDs)
			println("添加多分类筛选：category_id IN", validCategoryIDs)
		} else {
			println("没有有效的分类ID，不进行分类筛选")
		}
	} else {
		println("没有传递category_id参数，显示全部商品")
	}

	// 3. 只显示上架商品
	db = db.Where("status = 1")
	println("添加状态筛选：status = 1")

	// 4. 执行查询
	if err := db.Find(&products).Error; err != nil {
		println("查询失败:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品查询失败"})
		return
	}

	println("查询成功，返回商品数量:", len(products))
	println("=====================================")

	c.JSON(http.StatusOK, products)
}

// 获取商品详情
func GetProductDetail(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// 新增商品
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": product})
}

// 获取所有父分类
func GetParentCategories(c *gin.Context) {
	var list []models.Category
	config.DB.Where("parent_id = 0 AND status = 1").Find(&list)
	c.JSON(http.StatusOK, list)
}

// 获取子分类
func GetChildCategories(c *gin.Context) {
	parentID := c.Query("parent_id")
	var list []models.Category
	config.DB.Where("parent_id = ? AND status = 1", parentID).Find(&list)
	c.JSON(http.StatusOK, list)
}
