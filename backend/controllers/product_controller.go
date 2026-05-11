package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取商品列表
func GetProductList(c *gin.Context) {
	var products []models.Product
	keyword := c.Query("keyword")
	categoryID := c.Query("category_id")
	db := config.DB
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}
	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// 获取商品详情
func GetProductDetail(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// 新增商品（管理端）
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": product})
}
