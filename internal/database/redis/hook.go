package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/logger"
	"time"
)

// Hook 实现 redis.Hook 接口，用于日志记录
type Hook struct{}

func (h *Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	stats := CallStatsFromCtx(ctx)
	stats.Start = time.Now()
	ctx = WithCallStatsCtx(ctx, stats)
	return ctx, nil
}

func (h *Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	var (
		stats   = CallStatsFromCtx(ctx)
		latency = time.Since(stats.Start)
		logMap  = map[string]interface{}{
			logger.KAttempt:  stats.Attempt,
			logger.KRetries:  stats.Retries,
			"name":           cmd.Name(),
			"params":         cmd.Args(),
			logger.KProcTime: latency / time.Millisecond,
		}
		err = cmd.Err()
	)
	if err != nil {
		var (
			errType = errorx.ErrTypeSys
			isSuc   = errorx.Failed
			code    = errorx.ErrDefault.Code
		)
		if errors.Is(err, redis.Nil) {
			errType = errorx.ErrTypeBiz
			isSuc = errorx.Success
			code = errorx.ErrNotFound.Code
		}
		err = errorx.New(errorx.NewQuery{
			ErrMeta: errorx.ErrMeta{
				ServiceType: errorx.ServiceTypeBasic,
				Service:     errorx.ServiceRedis,
				ErrType:     errType,
				IsSuccess:   isSuc,
			},
			Code: code,
			Err:  err,
		})
	}

	mErr := errorx.Get(err, false)
	if mErr != nil {
		logMap[logger.KErr] = mErr
		logMap[logger.KErrMsg] = mErr.FinalMsg
		logMap[logger.KCode] = mErr.FinalCode
	}
	if mErr != nil && mErr.ErrType == errorx.ErrTypeSys {
		logger.Warn(ctx, logger.TagRedisSuccess, logMap)
	} else {
		logger.Info(ctx, logger.TagRedisSuccess, logMap)
	}
	return err
}

func (h *Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	fmt.Printf("[BeforeProcessPipeline] Pipeline with %d cmds\n", len(cmds))
	return ctx, nil
}

func (h *Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	fmt.Printf("[AfterProcessPipeline] Pipeline executed\n")
	return nil
}
