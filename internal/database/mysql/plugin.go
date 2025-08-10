package mysql

import (
	"time"

	"github.com/imattdu/go-web-v2/internal/common/cctx"
	"github.com/imattdu/go-web-v2/internal/common/logger"

	"gorm.io/gorm"
)

type LogPlugin struct{}

func (l LogPlugin) Name() string {
	return "logPlugin"
}

func (l LogPlugin) Initialize(db *gorm.DB) (err error) {
	beforeFuncName := "logPluginBefore"
	_ = db.Callback().Create().Before("gorm:before_create").Register(beforeFuncName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(beforeFuncName, before)
	_ = db.Callback().Update().Before("gorm:before_update").Register(beforeFuncName, before)
	_ = db.Callback().Query().Before("gorm:before_query").Register(beforeFuncName, before)
	_ = db.Callback().Raw().Before("gorm:before_raw").Register(beforeFuncName, before)
	_ = db.Callback().Row().Before("gorm:before_row").Register(beforeFuncName, before)

	afterFuncName := "logPluginAfter"
	_ = db.Callback().Create().After("gorm:after_create").Register(afterFuncName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(afterFuncName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(afterFuncName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(afterFuncName, after)
	_ = db.Callback().Raw().After("gorm:after_raw").Register(afterFuncName, after)
	_ = db.Callback().Row().After("gorm:after_row").Register(afterFuncName, after)
	return
}

func before(db *gorm.DB) {
	ctx := db.Statement.Context
	stats := CallStatsFromCtx(ctx)
	stats.Start = time.Now()
	ctx = WithCallStatsCtx(ctx, stats)

	trace := cctx.TraceFromCtxOrNew(ctx, nil).Copy()
	trace.UpdateParentSpanID()
	ctx = cctx.WithTraceCtx(ctx, trace)
	db.Statement.Context = ctx
}

func after(db *gorm.DB) {
	var (
		ctx      = db.Statement.Context
		stats    = CallStatsFromCtx(ctx)
		procTime = time.Since(stats.Start)
	)

	var (
		sql    = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars)
		err    = db.Error
		logMap = map[string]interface{}{
			"params":    stats.Params,
			"sql":       sql,
			"proc_time": procTime.Milliseconds(),
		}
	)
	if err != nil {
		logMap["err"] = err.Error()
		logger.Warn(ctx, logger.TagMysqlFailure, logMap)
		return
	}
	logger.Info(ctx, logger.TagMysqlSuccess, logMap)
}
