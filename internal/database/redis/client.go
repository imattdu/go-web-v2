package redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	GlobalRedisClient *redis.Client
	GlobalRetryClient *retryClient
)

func Init() {
	GlobalRedisClient = newRedisClient()
	GlobalRetryClient = newRetryClient(retryConf{
		retries: 3,
		delay:   time.Millisecond * 100,
	})
}

func newRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
	})
	// 注册 Hook
	rdb.AddHook(&Hook{})
	return rdb
}
