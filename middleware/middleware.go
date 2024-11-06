package middleware

import (
	"errors"
	"github.com/GitEval/GitEval-Backend/api/response"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

var ProviderSet = wire.NewSet(NewMiddleware, NewJWTClient)

type ParTokener interface {
	ParseToken(tokenString string) (int64, error)
}
type Middleware struct {
	jwt ParTokener
}

func NewMiddleware(jwt ParTokener) *Middleware {
	return &Middleware{jwt}
}

// AuthMiddleware 从请求头中获取认证信息并解析出 user_id
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Err{Err: errors.New("Authorization header is empty.")})
			c.Abort()
			return
		}

		//解析jwt
		userID, err := m.jwt.ParseToken(authHeader)
		if err != nil || userID == 0 {
			c.JSON(http.StatusUnauthorized, response.Err{Err: err})
			c.Abort()
			return
		}

		// 将 user_id 存储到上下文中
		c.Set("user_id", userID)

		// 继续处理请求
		c.Next()
	}
}
