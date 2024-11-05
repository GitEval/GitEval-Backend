package middleware

import (
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type JWTClient struct {
	cfg *conf.JWTConfig
}

func NewJWTClient(config *conf.JWTConfig) *JWTClient {
	return &JWTClient{cfg: config}
}

// GenerateToken 生成 ParTokener token
func (c *JWTClient) GenerateToken(userID int64) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建 token
	claims := &jwt.StandardClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签署 token
	return token.SignedString([]byte(c.cfg.SecretKey))
}

// ParseToken 解析 ParTokener token 并返回 userID
func (c *JWTClient) ParseToken(tokenString string) (int64, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorMalformed)
		}
		return []byte(c.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}
	// 转换为 int64
	userId, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}
	return userId, nil // 返回 userID
}
