package redis

import (
	"context"
	"math/rand"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/util"
)

// KVEntry 支持 key 前缀、key 主体、value 主体及带抖动的 TTL
type KVEntry struct {
	KeyPre  string      `json:"key_pre"`  // Key 前缀
	KeyBody interface{} `json:"key_body"` // Key 主体
	VBody   interface{} `json:"v_body"`   // Value 主体，支持任意类型

	UseKeyBody    bool `json:"use_key_body"`
	EncryptKey    bool `json:"encrypt_key"`
	CompressValue bool `json:"compress_value"`

	BaseTTL   time.Duration `json:"base_ttl"`   // 基础 TTL
	MaxJitter time.Duration `json:"max_jitter"` // 最大随机抖动 TTL
}

func (e *KVEntry) Key(ctx context.Context) (string, error) {
	var k = e.KeyPre
	if !e.UseKeyBody {
		return k, nil
	}
	kBodyStr, err := util.Marshal(ctx, e.KeyBody)
	if err != nil {
		return "", err
	}
	if e.EncryptKey {
		kBodyStr = util.Sha256Hex([]byte(kBodyStr))
	}
	return k + kBodyStr, nil
}

func (e *KVEntry) Value(ctx context.Context) (string, error) {
	vBodyStr, err := util.Marshal(ctx, e.VBody)
	if err != nil {
		return "", err
	}
	if e.CompressValue {
		vBodyBytes, err := util.Compress([]byte(vBodyStr))
		if err != nil {
			return "", err
		}
		vBodyStr = string(vBodyBytes)
	}
	return vBodyStr, nil
}

func (e *KVEntry) ValueRaw(ctx context.Context, value string) (string, error) {
	if e.CompressValue {
		valueBytes, err := util.Compress([]byte(value))
		if err != nil {
			return "", err
		}
		value = string(valueBytes)
	}
	return value, nil
}

func (e *KVEntry) TTL() time.Duration {
	if e.MaxJitter <= 0 {
		return e.BaseTTL
	}
	jitter := time.Duration(rand.Int63n(int64(e.MaxJitter)))
	return e.BaseTTL + jitter
}
