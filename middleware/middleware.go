package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware 从请求头中获取认证信息并解析出 user_id
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		//解析jwt
		userID, err := ParseToken(authHeader)
		if err != nil || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}

		// 将 user_id 存储到上下文中
		c.Set("user_id", userID)

		// 继续处理请求
		c.Next()
	}
}
