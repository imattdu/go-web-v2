package logger

import (
	"context"
	"log/slog"
)

func Info(ctx context.Context, tag string, log interface{}) {
	Logger.InfoContext(context.WithoutCancel(ctx), "", slog.Any(logK, Log{
		caller:  getCaller(),
		tag:     tag,
		message: log,
	}))
}

func Warn(ctx context.Context, tag string, log interface{}) {
	Logger.WarnContext(context.WithoutCancel(ctx), "", slog.Any(logK, Log{
		caller:  getCaller(),
		tag:     tag,
		message: log,
	}))
}

func Error(ctx context.Context, tag string, log interface{}) {
	Logger.ErrorContext(context.WithoutCancel(ctx), "", slog.Any(logK, Log{
		caller:  getCaller(),
		tag:     tag,
		message: log,
	}))
}
