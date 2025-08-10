package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type caller struct {
	file     string
	line     int
	funcName string
}

type Log struct {
	caller  caller
	tag     string
	message interface{}
}

type logEntry struct {
	ctx    context.Context
	record slog.Record
	done   chan error
}

type LogConfig struct {
	LogDir         string     // 日志目录
	Level          slog.Level // 最小日志级别
	AppName        string     // 应用名
	ConsoleColored bool       // 控制台是否彩色输出
	MaxFileSizeMB  int        // 超过多少 MB 分割，0 表示按小时切割
	MaxBackups     int        // 最多保留文件个数
	Timeout        time.Duration
	ConsoleEnabled bool
}

type zeroHandler struct {
	slog.Handler
	cfg          LogConfig
	curFile      *os.File
	curWarnFile  *os.File
	curSize      int64
	curWarnSize  int64
	curIndex     int
	curWarnIndex int
	curTime      time.Time
	entries      chan logEntry
}
