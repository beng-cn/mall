package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 新增商品
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		fmt.Printf("参数绑定失败：%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	if strings.TrimSpace(product.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品名称不能为空"})
		return
	}
	if product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品价格必须大于0"})
		return
	}
	if product.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "库存数量不能为负数"})
		return
	}
	if product.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择商品分类"})
		return
	}
	if product.Status != 1 && product.Status != 0 {
		product.Status = 1
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品创建失败：" + err.Error()})
		return
	}

	// ✅ 清理缓存
	ClearAllProductListCache()

	c.JSON(http.StatusOK, gin.H{"message": "商品创建成功", "data": product})
}

// 添加分类
func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	if strings.TrimSpace(category.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类名称不能为空"})
		return
	}
	if category.Status != 1 && category.Status != 0 {
		category.Status = 1
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分类创建失败：" + err.Error()})
		return
	}

	// ✅ 清理缓存
	ClearAllCategoryCache()

	c.JSON(http.StatusOK, gin.H{"message": "分类创建成功", "data": category})
}

// 编辑分类
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	if strings.TrimSpace(category.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类名称不能为空"})
		return
	}

	if err := config.DB.Model(&models.Category{}).Where("id = ?", id).Update("name", category.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分类更新失败：" + err.Error()})
		return
	}

	// ✅ 清理缓存
	ClearAllCategoryCache()

	c.JSON(http.StatusOK, gin.H{"message": "分类更新成功"})
}

// 删除分类
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	var productCount int64
	config.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
	if productCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下存在商品，无法删除"})
		return
	}

	var childCount int64
	config.DB.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下存在子分类，无法删除"})
		return
	}

	if err := config.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分类删除失败：" + err.Error()})
		return
	}

	// ✅ 清理缓存
	ClearAllCategoryCache()

	c.JSON(http.StatusOK, gin.H{"message": "分类删除成功"})
}

// 上传图片
func UploadImage(c *gin.Context) {
	c.Request.ParseMultipartForm(10 << 20)
	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取上传文件失败：" + err.Error()})
		return
	}
	defer file.Close()

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}
	contentType := handler.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持上传 jpg、png、gif"})
		return
	}

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	dst, _ := os.Create(filePath)
	defer dst.Close()
	dst.ReadFrom(file)

	imageURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": imageURL})
}

// 更新商品
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	if strings.TrimSpace(updateData.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品名称不能为空"})
		return
	}
	if updateData.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品价格必须大于0"})
		return
	}
	if updateData.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "库存数量不能为负数"})
		return
	}
	if updateData.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择商品分类"})
		return
	}

	config.DB.Model(&product).Updates(updateData)

	// ✅ 清理缓存
	ClearSingleProductCache(id)
	ClearAllProductListCache()

	c.JSON(http.StatusOK, gin.H{"message": "商品更新成功"})
}

// 删除商品
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	config.DB.Delete(&product)

	// ✅ 清理缓存
	ClearSingleProductCache(id)
	ClearAllProductListCache()

	c.JSON(http.StatusOK, gin.H{"message": "商品删除成功"})
}
