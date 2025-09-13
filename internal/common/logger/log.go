package logger

import (
	"context"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"log/slog"
)

func Info(ctx context.Context, tag string, log interface{}) {
	Logger.InfoContext(cctx.CloneWithoutDeadline(ctx), "", slog.Any(logK, Log{
		caller:  getCaller(),
		tag:     tag,
		message: log,
	}))
}

func Warn(ctx context.Context, tag string, log interface{}) {
	Logger.WarnContext(cctx.CloneWithoutDeadline(ctx), "", slog.Any(logK, Log{
		caller:  getCaller(),
		tag:     tag,
		message: log,
	}))
}

func Error(ctx context.Context, tag string, log interface{}) {
	Logger.ErrorContext(cctx.CloneWithoutDeadline(ctx), "", slog.Any(logK, Log{
		caller:  getCaller(),
		tag:     tag,
		message: log,
	}))
}
