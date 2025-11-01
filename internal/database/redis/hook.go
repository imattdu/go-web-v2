package redis

import (
	"context"
	"errors"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/logger"

	"github.com/go-redis/redis/v8"
)

// Hook 实现 redis.Hook 接口，用于日志记录
type Hook struct{}

func (h *Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	stats, ok := cctx.GetAs[*CallStats](ctx, cctx.RedisCallStatsKey)
	if !ok {
		return ctx, nil
	}
	stats.Start = time.Now()
	return ctx, nil
}

func (h *Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	stats, ok := cctx.GetAs[*CallStats](ctx, cctx.RedisCallStatsKey)
	if !ok {
		return nil
	}

	var (
		latency = time.Since(stats.Start)
		logMap  = map[string]interface{}{
			logger.KAttempt:    stats.Attempt,
			logger.KRetries:    stats.Retries,
			logger.KCalleeFunc: cmd.Name(),
			logger.KCmd:        cmd.Args(),
			logger.KProcTime:   latency / time.Millisecond,
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
		err = errorx.New(errorx.ErrOptions{
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
	isFinal := stats.Attempt == stats.Retries
	if !isFinal && !shouldRetry(err) {
		isFinal = true
	}
	logMap[logger.KIsFinal] = isFinal

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
	return ctx, nil
}

func (h *Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
