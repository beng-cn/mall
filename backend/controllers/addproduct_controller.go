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

// 新增商品（完善版，含完整参数校验+缓存清理）
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		fmt.Printf("参数绑定失败：%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	// 核心参数校验
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
	// 上架状态默认1（上架）
	if product.Status != 1 && product.Status != 0 {
		product.Status = 1
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品创建失败：" + err.Error()})
		return
	}

	// 清理所有商品列表缓存
	keys, _ := config.RDB.Keys(config.Ctx, "product:list:*").Result()
	if len(keys) > 0 {
		config.RDB.Del(config.Ctx, keys...)
	}

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
	// 状态默认1（启用）
	if category.Status != 1 && category.Status != 0 {
		category.Status = 1
	}

	// 写入数据库
	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分类创建失败：" + err.Error()})
		return
	}

	// 清理对应缓存
	if category.ParentID == 0 {
		config.RDB.Del(config.Ctx, "category:parent")
	} else {
		cacheKey := fmt.Sprintf("category:child:%d", category.ParentID)
		config.RDB.Del(config.Ctx, cacheKey)
	}

	c.JSON(http.StatusOK, gin.H{"message": "分类创建成功", "data": category})
}

// 编辑分类名称
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

	// 更新数据库
	if err := config.DB.Model(&models.Category{}).Where("id = ?", id).Update("name", category.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分类更新失败：" + err.Error()})
		return
	}

	// 清理所有分类缓存
	config.RDB.Del(config.Ctx, "category:parent")
	childKeys, _ := config.RDB.Keys(config.Ctx, "category:child:*").Result()
	if len(childKeys) > 0 {
		config.RDB.Del(config.Ctx, childKeys...)
	}

	c.JSON(http.StatusOK, gin.H{"message": "分类更新成功"})
}

// 删除分类（安全校验：有商品/子分类则禁止删除）
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// 1. 检查该分类下是否有商品
	var productCount int64
	config.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
	if productCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下存在商品，无法删除"})
		return
	}

	// 2. 检查该分类下是否有子分类
	var childCount int64
	config.DB.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下存在子分类，无法删除"})
		return
	}

	// 3. 执行删除
	if err := config.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分类删除失败：" + err.Error()})
		return
	}

	// 4. 清理缓存
	config.RDB.Del(config.Ctx, "category:parent")
	childKeys, _ := config.RDB.Keys(config.Ctx, "category:child:*").Result()
	if len(childKeys) > 0 {
		config.RDB.Del(config.Ctx, childKeys...)
	}

	c.JSON(http.StatusOK, gin.H{"message": "分类删除成功"})
}

// 商品图片上传接口（支持jpg/png/gif，最大10MB）
func UploadImage(c *gin.Context) {
	// 限制单文件大小为10MB
	c.Request.ParseMultipartForm(10 << 20)

	// 获取上传的文件
	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取上传文件失败：" + err.Error()})
		return
	}
	defer file.Close()

	// 校验文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}
	contentType := handler.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持上传 jpg、png、gif 格式的图片"})
		return
	}

	// 创建上传目录（不存在则自动创建）
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
			return
		}
	}

	// 生成唯一文件名
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件到服务器
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败：" + err.Error()})
		return
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败：" + err.Error()})
		return
	}

	// 返回可访问的图片URL
	imageURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": imageURL})
}

// 更新商品（编辑已有商品）
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// 先查询商品是否存在
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	// 绑定更新数据
	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	// 核心参数校验
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

	// 更新数据库
	if err := config.DB.Model(&product).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品更新失败：" + err.Error()})
		return
	}

	// 清理缓存：商品列表 + 该商品详情
	keys, _ := config.RDB.Keys(config.Ctx, "product:list:*").Result()
	if len(keys) > 0 {
		config.RDB.Del(config.Ctx, keys...)
	}
	config.RDB.Del(config.Ctx, "product:item:"+id)

	c.JSON(http.StatusOK, gin.H{"message": "商品更新成功", "data": product})
}

// 删除商品
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// 先查询商品是否存在
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	// 执行软删除（GORM默认软删除，保留数据在数据库）
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "商品删除失败：" + err.Error()})
		return
	}

	// 清理缓存：所有商品列表 + 该商品详情
	keys, _ := config.RDB.Keys(config.Ctx, "product:list:*").Result()
	if len(keys) > 0 {
		config.RDB.Del(config.Ctx, keys...)
	}
	config.RDB.Del(config.Ctx, "product:item:"+id)

	c.JSON(http.StatusOK, gin.H{"message": "商品删除成功"})
}
