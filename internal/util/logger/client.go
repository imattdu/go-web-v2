package logger

import (
	"log/slog"
	"time"

	"github.com/imattdu/go-web-v2/internal/config"
)

var Logger *slog.Logger

func Init() {
	var (
		logConf = config.GlobalConf.Log
		conf    = LogConfig{
			LogDir:         logConf.Path,
			Level:          slog.LevelInfo,
			AppName:        logConf.FileName,
			ConsoleColored: true,
			MaxBackups:     logConf.MaxBackups,
			Timeout:        time.Millisecond * logConf.Timeout,
		}
	)
	if logConf.MaxBackups <= 4 {
		logConf.MaxBackups = 4
	}
	if logConf.Timeout == 0 {
		logConf.Timeout = 2000 * time.Millisecond
	}
	handler := NewZeroHandler(conf)
	Logger = slog.New(handler)
}
