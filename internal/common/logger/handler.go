package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	trace2 "github.com/imattdu/go-web-v2/internal/common/trace"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewZeroHandler(cfg LogConfig) slog.Handler {
	h := &zeroHandler{
		cfg:     cfg,
		curTime: time.Now().Truncate(time.Hour),
		entries: make(chan logEntry, 1000), // 队列大小可调
	}

	h.curFile = h.openNewFile(false)
	h.curWarnFile = h.openNewFile(true)
	h.cleanupOldFiles()
	h.curSize = 0
	h.curWarnSize = 0
	// 启动写日志协程
	go h.writeLoop()
	return h
}

func (h *zeroHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.cfg.Level
}

func (h *zeroHandler) Handle(ctx context.Context, r slog.Record) error {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.Timeout)
	defer cancel() // 必须取消，以防资源泄露

	// 复制 Record，防止异步读取问题
	rr := r.Clone()
	le := logEntry{
		ctx:    ctx,
		record: rr,
		done:   make(chan error, 1),
	}

	// 发送写入请求（阻塞等待写完）
	select {
	case h.entries <- le:
		select {
		case err := <-le.done:
			return err
		case <-ctx.Done():
			log.Println("写日志超时:", ctx.Err())
			return ctx.Err()
		}
	case <-ctx.Done():
		log.Println("发送日志超时", ctx.Err())
		return ctx.Err()
	}
}

func (h *zeroHandler) writeLoop() {
	for le := range h.entries {
		err := h.writeRecord(le.ctx, le.record)
		le.done <- err
	}
}

func (h *zeroHandler) writeRecord(ctx context.Context, r slog.Record) error {
	now := time.Now()
	if h.shouldRotate(now, false) {
		h.rotate(now, false)
	}
	if h.shouldRotate(now, true) {
		h.rotate(now, true)
	}

	var sb strings.Builder
	parts := []string{
		r.Level.String(),
		r.Time.Format(time.RFC3339),
	}
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != logK {
			return true
		}

		l, ok := a.Value.Any().(Log)
		if !ok {
			return false
		}
		parts = append(parts, l.caller.funcName+"@"+l.caller.file+":"+strconv.Itoa(l.caller.line))
		for _, part := range parts {
			sb.WriteString("[")
			sb.WriteString(part)
			sb.WriteString("]")
		}

		sb.WriteString(" ")
		sb.WriteString(l.tag)
		sb.WriteString("||")
		sb.WriteString(cctx.TraceFromCtxOrNew(ctx, func() *trace2.Trace {
			return trace2.New(&http.Request{URL: &url.URL{}})
		}).String())

		msgMap, ok := l.message.(map[string]interface{})
		if ok {
			for k, v := range msgMap {
				sb.WriteString(k)
				sb.WriteString("=")
				msg, _ := json.Marshal(v)
				sb.WriteString(string(msg))
				sb.WriteString("||")
			}
		} else {
			msg, _ := json.Marshal(l.message)
			sb.WriteString("msg=")
			sb.WriteString(string(msg))
		}

		return false
	})
	line := sb.String() + "\n"
	if h.cfg.ConsoleEnabled {
		fmt.Printf(line)
	}

	var (
		n   int
		err error
	)
	if r.Level >= slog.LevelWarn {
		n, err = h.curWarnFile.WriteString(line)
		if err != nil {
			//
		}
		h.curWarnSize += int64(n)
	} else {
		n, err = h.curFile.WriteString(line)
		h.curSize += int64(n)
	}
	return err
}

func (h *zeroHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *zeroHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *zeroHandler) shouldRotate(now time.Time, isWarnFile bool) bool {
	curSize := h.curSize
	if isWarnFile {
		curSize = h.curWarnSize
	}
	if h.cfg.MaxFileSizeMB > 0 && curSize >= int64(h.cfg.MaxFileSizeMB)<<20 {
		return true
	}
	if now.Truncate(time.Hour) != h.curTime {
		return true
	}
	return false
}

func (h *zeroHandler) rotate(now time.Time, isWarnFile bool) {
	if isWarnFile {
		if h.curWarnFile != nil {
			_ = h.curWarnFile.Close()
		}
		// 时间切割时，告警文件序号重置
		if now.Truncate(time.Hour) != h.curTime {
			h.curWarnIndex = 1
		} else {
			h.curWarnIndex++
		}
		h.curWarnFile = h.openNewFile(true)
		h.curWarnSize = 0
		h.curTime = now.Truncate(time.Hour)
	} else {
		if h.curFile != nil {
			_ = h.curFile.Close()
		}
		if now.Truncate(time.Hour) != h.curTime {
			h.curIndex = 1
			h.curTime = now.Truncate(time.Hour)
		} else {
			h.curIndex++
		}
		h.curFile = h.openNewFile(false)
		h.curSize = 0
		h.curTime = now.Truncate(time.Hour)
	}
	h.cleanupOldFiles()
}

func (h *zeroHandler) openNewFile(isWarnFile bool) *os.File {
	if err := os.MkdirAll(h.cfg.LogDir, 0755); err != nil {
		log.Println(err.Error())
		return nil
	}
	var (
		suffix = ".log"
		link   = h.cfg.AppName + ".log"
	)
	if isWarnFile {
		suffix += ".wf"
		link += ".wf"
	}
	suffix += "." + h.curTime.Format("2006010215")
	if h.cfg.MaxFileSizeMB > 0 {
		if isWarnFile {
			suffix = fmt.Sprintf("%s.%03d", suffix, h.curWarnIndex)
		} else {
			suffix = fmt.Sprintf("%s.%03d", suffix, h.curIndex)
		}
	}

	var (
		file = h.cfg.AppName + suffix
		path = filepath.Join(h.cfg.LogDir, file)
	)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	linkPath := filepath.Join(h.cfg.LogDir, link)
	if err = os.Remove(linkPath); err != nil {
		fmt.Println(err.Error())
	}
	if err = os.Symlink(file, linkPath); err != nil {
		fmt.Println(err.Error())
	}
	return f
}

func (h *zeroHandler) cleanupOldFiles() {
	dirEntries, err := os.ReadDir(h.cfg.LogDir)
	if err != nil || h.cfg.MaxBackups <= 0 {
		return
	}

	var logFiles []string
	prefix := h.cfg.AppName
	for _, e := range dirEntries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), prefix) {
			logFiles = append(logFiles, e.Name())
		}
	}
	if len(logFiles) <= h.cfg.MaxBackups {
		return
	}
	sort.Slice(logFiles, func(i, j int) bool {
		var (
			nextNodes = strings.Split(logFiles[i], ".")
			nodes     = strings.Split(logFiles[j], ".")
		)
		nextDt, err := time.Parse("2006010215", nextNodes[len(nextNodes)-1])
		if err != nil {
			return true
		}
		dt, err := time.Parse("2006010215", nodes[len(nodes)-1])
		if err != nil {
			return false
		}
		return nextDt.After(dt)
	})
	remove := logFiles[h.cfg.MaxBackups:]
	for _, f := range remove {
		_ = os.Remove(filepath.Join(h.cfg.LogDir, f))
	}
}

func (h *zeroHandler) colorLine(r slog.Record) string {
	color := ""
	emoji := ""
	switch r.Level {
	case slog.LevelDebug:
		color = "\033[36m" // Cyan
		emoji = "🐛"
	case slog.LevelInfo:
		color = "\033[32m"
		emoji = "ℹ️"
	case slog.LevelWarn:
		color = "\033[33m"
		emoji = "⚠️"
	case slog.LevelError:
		color = "\033[31m"
		emoji = "❌"
	default:
		color = "\033[0m"
	}
	line := fmt.Sprintf("%s[%s %s] %s\033[0m\n", color, emoji, r.Level, r.Message)
	return line
}
