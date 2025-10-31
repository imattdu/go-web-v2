package redis

import (
	"time"
)

type CallStats struct {
	Attempt int
	Retries int
	Start   time.Time
}
