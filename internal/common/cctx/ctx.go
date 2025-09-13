package cctx

import (
	"context"
	"sync"
)

// ---------- 类型定义 ----------
type ctxKey string

var dataKey ctxKey = "ctx_data"

// Metadata 类型
type metadata struct {
	mu   sync.RWMutex
	data map[string]any
}

func New(parent context.Context, data map[string]any) context.Context {
	md := &metadata{
		data: deepCopyMap(data),
	}
	return context.WithValue(parent, dataKey, md)
}

func fromContext(ctx context.Context) *metadata {
	if v, ok := ctx.Value(dataKey).(*metadata); ok {
		return v
	}
	// 返回空的只读 Metadata
	return &metadata{
		data: make(map[string]any),
	}
}

func (md *metadata) get(key string) (any, bool) {
	md.mu.RLock()
	defer md.mu.RUnlock()
	val, ok := md.data[key]
	return val, ok
}

func (md *metadata) set(key string, val any) {
	md.mu.Lock()
	defer md.mu.Unlock()
	md.data[key] = val
}

func (md *metadata) all() map[string]any {
	md.mu.RLock()
	defer md.mu.RUnlock()
	return deepCopyMap(md.data)
}

func Get(ctx context.Context, key string) (any, bool) {
	md := fromContext(ctx)
	return md.get(key)
}

func Set(ctx context.Context, key string, val any) {
	md := fromContext(ctx)
	md.set(key, val)
}

func All(ctx context.Context) map[string]any {
	md := fromContext(ctx)
	return md.all()
}

// ---------- 泛型接口 ----------

// GetAs 从 ctx 获取指定类型的 value
func GetAs[T any](ctx context.Context, key string) (T, bool) {
	var zero T
	md := fromContext(ctx)
	val, ok := md.get(key)
	if !ok {
		return zero, false
	}
	tVal, ok := val.(T)
	if !ok {
		return zero, false
	}
	return tVal, true
}

// SetAs 将任意类型 value 存入 ctx
func SetAs[T any](ctx context.Context, key string, val T) {
	md := fromContext(ctx)
	md.set(key, val)
}

// AllAs 获取 ctx 内所有能转换为 T 的键值对
func AllAs[T any](ctx context.Context) map[string]T {
	md := fromContext(ctx)
	result := make(map[string]T)
	for k, v := range md.all() {
		if val, ok := v.(T); ok {
			result[k] = val
		}
	}
	return result
}

func Clone(ctx context.Context) (context.Context, context.CancelFunc) {
	var newCtx context.Context
	var cancel context.CancelFunc

	// 继承 deadline / cancel
	if deadline, ok := ctx.Deadline(); ok {
		newCtx, cancel = context.WithDeadline(context.Background(), deadline)
	} else {
		newCtx, cancel = context.WithCancel(context.Background())
	}

	orig := fromContext(ctx)
	dataCopy := orig.all() // 深拷贝 map
	newCtx = New(newCtx, dataCopy)
	return newCtx, cancel
}

func CloneWithoutDeadline(ctx context.Context) context.Context {
	orig := fromContext(ctx)
	dataCopy := orig.all()
	return New(context.Background(), dataCopy)
}

// ---------- 辅助函数 ----------
func deepCopyMap(orig map[string]any) map[string]any {
	newMap := make(map[string]any, len(orig))
	for k, v := range orig {
		newMap[k] = v
	}
	return newMap
}
