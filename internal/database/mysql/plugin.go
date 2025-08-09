package mysql

import (
	"github.com/imattdu/go-web-v2/internal/common/cctx"
	logger2 "github.com/imattdu/go-web-v2/internal/common/util/logger"
	"gorm.io/gorm"
	"time"
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
	d := cctx.MysqlFromCtx(ctx)
	d.Start = time.Now()
	d.Query = nil
	ctx = cctx.WithMysqlCtx(ctx, d)

	trace := cctx.TraceFromCtxOrNew(ctx, nil).Copy()
	trace.UpdateParentSpanID()
	ctx = cctx.WithTraceCtx(ctx, trace)
	db.Statement.Context = ctx
}

func after(db *gorm.DB) {
	var (
		ctx      = db.Statement.Context
		d        = cctx.MysqlFromCtx(ctx)
		procTime = time.Since(d.Start)
	)

	var (
		sql    = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars)
		err    = db.Error
		logMap = map[string]interface{}{
			"query":     d.Query,
			"sql":       sql,
			"proc_time": procTime.Milliseconds(),
		}
	)
	if err != nil {
		logMap["err"] = err.Error()
		logger2.Warn(ctx, logger2.TagMysqlFailure, logMap)
		return
	}
	logger2.Info(ctx, logger2.TagMysqlSuccess, logMap)
}
