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

	db := config.DB

	// 关键词模糊搜索
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 多分类筛选
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
			}
		}

		if len(validCategoryIDs) > 0 {
			db = db.Where("category_id IN (?)", validCategoryIDs)
		}
	}

	// 只显示上架商品
	db = db.Where("status = 1")

	// 查询
	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品查询失败"})
		return
	}

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
