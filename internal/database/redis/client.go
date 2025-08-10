package redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	GlobalRedisClient *redis.Client
	GlobalRetryClient RetryClient
)

func Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		//DB:       0,  // use default DB
	})
	// 注册 Hook
	rdb.AddHook(&Hook{})

	GlobalRedisClient = rdb

	GlobalRetryClient = RetryClient{
		retries: 3,
		delay:   100 * time.Millisecond,
	}
}
