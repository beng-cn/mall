package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 注册
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 校验用户名和密码
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}
	if len(user.Username) < 3 || len(user.Username) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名长度需在3-20位之间"})
		return
	}
	if len(user.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度不能少于6位"})
		return
	}

	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	user.Password = string(hash)

	// 写入数据库
	if err := config.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败：" + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// 登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 校验用户状态
	if user.Status != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户已被禁用"})
		return
	}

	// 解密比对密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.RoleID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 登录成功
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user":    user,
	})
}

// 修改用户信息
func UpdateUserInfo(c *gin.Context) {
	loginUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		// 不允许修改ID、RoleID、Status等敏感字段
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 强制更新当前登录用户的信息
	if err := config.DB.Model(&models.User{}).Where("id = ?", loginUserID.(uint)).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// 管理员获取用户列表
func ListUsers(c *gin.Context) {
	// 获取分页和搜索参数
	pageNum, _ := strconv.Atoi(c.Query("page_num"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	keyword := c.Query("keyword")

	// 分页默认值
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建查询条件
	var users []models.User
	var total int64
	tx := config.DB.Model(&models.User{})

	// 搜索关键词
	if keyword != "" {
		tx = tx.Where("username LIKE ?", "%"+keyword+"%")
	}

	// 统计总数
	tx.Count(&total)

	// 分页查询
	offset := (pageNum - 1) * pageSize
	if err := tx.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户列表失败",
		})
		return
	}

	// 隐藏密码
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  users,
		"total": total,
	})
}

// 管理员更新用户状态
func UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误：" + err.Error()})
		return
	}

	// 更新状态
	if err := config.DB.Model(&models.User{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// 管理员删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 禁止删除管理员
	var user models.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	if user.RoleID == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "禁止删除管理员用户"})
		return
	}

	// 软删除
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除用户成功"})
}
