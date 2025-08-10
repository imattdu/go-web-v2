package redis

import (
	"context"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/util"
)

func Set(ctx context.Context, e *KVEntry) error {
	k, err := e.Key(ctx)
	if err != nil {
		return err
	}
	v, err := e.Value(ctx)
	if err != nil {
		return err
	}

	return GlobalRetryClient.RetryCmd(ctx, func(mewCtx context.Context) error {
		return GlobalRedisClient.Set(mewCtx, k, v, e.TTL()).Err()
	})
}

func Get(ctx context.Context, e *KVEntry) (string, error) {
	k, err := e.Key(ctx)
	if err != nil {
		return "", err
	}

	var value string
	err = GlobalRetryClient.RetryCmd(ctx, func(mewCtx context.Context) error {
		curValue, err := GlobalRedisClient.Get(mewCtx, k).Result()
		value = curValue
		return err
	})
	_ = util.Unmarshal(ctx, value, &e.VBody)
	return value, err
}

func TTL(ctx context.Context, e *KVEntry) (time.Duration, error) {
	k, err := e.Key(ctx)
	if err != nil {
		return 0, err
	}

	var ttl time.Duration
	err = GlobalRetryClient.RetryCmd(ctx, func(mewCtx context.Context) error {
		ttl, err = GlobalRedisClient.TTL(mewCtx, k).Result()
		return err
	})
	return ttl, err
}
