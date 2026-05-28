package middlewares

import (
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 抽离核心登录校验逻辑，供两个权限中间件复用
func validateAuth(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录，请先登录"})
		c.Abort()
		return false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请求头格式错误"})
		c.Abort()
		return false
	}

	claims, err := utils.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "登录已过期，请重新登录"})
		c.Abort()
		return false
	}

	// 将用户信息存入上下文，供后续控制器使用
	c.Set("user_id", claims.UserID)
	c.Set("role_id", claims.RoleID)
	c.Set("username", claims.Username)
	return true
}

// 普通用户登录校验中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if validateAuth(c) {
			c.Next()
		}
	}
}

// 管理员权限校验中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !validateAuth(c) {
			return
		}

		// 校验管理员权限
		roleID, exists := c.Get("role_id")
		if !exists || roleID.(uint) != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "仅管理员可访问此接口"})
			c.Abort()
			return
		}

		c.Next()
	}
}
