package middleware

import (
	"errors"
	"github.com/GitEval/GitEval-Backend/api/response"
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/GitEval/GitEval-Backend/model/cache"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type JWTClient struct {
	cfg        *conf.JWTConfig
	redisCache *cache.RedisClient // 引入 Redis 缓存
}

func NewJWTClient(config *conf.JWTConfig, redisCache *cache.RedisClient) *JWTClient {
	return &JWTClient{cfg: config, redisCache: redisCache}
}

// GenerateToken 生成 ParTokener token
func (c *JWTClient) GenerateToken(userID int64) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(c.cfg.Timeout) * time.Minute)

	// 创建一个唯一的 jti
	jti := strconv.FormatInt(time.Now().UnixNano(), 10) // 使用当前时间的纳秒作为 jti，确保唯一性

	// 创建 token
	claims := &jwt.StandardClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: expirationTime.Unix(),
		Id:        jti, // 将 jti 作为 Id 字段
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

// BlackJWT 将 JWT 标记为无效，添加到黑名单
func (c *JWTClient) BlackJWT(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, response.Err{Err: errors.New("Authorization header is missing")})
		ctx.Abort()
		return
	}

	// 解析 JWT 获取用户 ID 和 jti
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, response.Err{Err: errors.New("Invalid token")})
		ctx.Abort()
		return
	}

	// 将 jti 添加到 Redis 黑名单，过期时间与 JWT 的有效时间相同
	expTime := time.Until(time.Unix(claims.ExpiresAt, 0))
	if err := c.redisCache.AddToBlacklist(claims.Id, int64(expTime.Seconds())); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: errors.New("blacklist set fail")})
		ctx.Abort()
		return
	}

	return
}
