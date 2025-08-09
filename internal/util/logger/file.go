package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getCaller() caller {
	pc, file, line, ok := runtime.Caller(2)
	file = trimFilePath(file)
	funcName := "unknown"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = trimFuncName(fn.Name())
		}
	}
	return caller{
		file:     file,
		line:     line,
		funcName: funcName,
	}
}

// 尝试从 go.mod 根目录起算相对路径
func trimFilePath(fullPath string) string {
	modRoot, err := findGoModRoot(fullPath)
	if err != nil {
		// fallback: 仅保留文件名
		_, short := filepath.Split(fullPath)
		return short
	}
	rel, err := filepath.Rel(modRoot, fullPath)
	if err != nil {
		return fullPath
	}
	return rel
}

// 向上查找 go.mod 根目录
func findGoModRoot(start string) (string, error) {
	dir := filepath.Dir(start)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found")
}

// 清理函数名（去掉包路径前缀）
func trimFuncName(name string) string {
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	// 去掉 package 路径，只保留函数名
	if idx := strings.Index(name, "."); idx >= 0 && idx+1 < len(name) {
		name = name[idx+1:]
	}
	return name
}
