package mysql

import (
	"time"
)

var callStatsKey = "mysqlCallStatsKey"

type CallStats struct {
	Params interface{}
	Start  time.Time
}
