package cache

import (
	"context"
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(conf *conf.CacheConf) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
	})
	return &RedisClient{client: client}
}

// AddToBlacklist 将 jti 添加到黑名单中
func (r *RedisClient) AddToBlacklist(jti string, expTime int64) error {
	return r.client.Set(context.Background(), jti, "blacklist", time.Duration(expTime)*time.Minute).Err()
}

// CheckBlacklist 检查 jti 是否在黑名单中
func (r *RedisClient) CheckBlacklist(jti string) (bool, error) {
	result, err := r.client.Exists(context.Background(), jti).Result()
	if err != nil {
		return false, err
	}
	// 如果 result > 0，说明该 token 被列入黑名单
	return result > 0, nil
}
