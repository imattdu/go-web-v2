package mysql

import (
	"errors"
	"github.com/imattdu/go-web-v2/internal/common/cctxv2"
	"github.com/imattdu/go-web-v2/internal/common/errorx"
	"github.com/imattdu/go-web-v2/internal/common/trace"
	"net/http"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/logger"

	"gorm.io/gorm"
)

type Plugin struct {
	callee string
}

func (l Plugin) Name() string {
	return "logPlugin"
}

func (l Plugin) Initialize(db *gorm.DB) (err error) {
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
	stats, ok := cctxv2.GetAs[*CallStats](ctx, cctxv2.MySQLCallStatsKey)
	if !ok {
		stats = &CallStats{}
	}
	stats.Start = time.Now()

	t, ok := cctxv2.GetAs[*trace.Trace](ctx, cctxv2.TraceKey)
	if !ok {
		t = trace.New(&http.Request{})
	}
	t = t.Copy()
	t.UpdateParentSpanID()
	ctx = cctxv2.With(ctx, cctxv2.TraceKey, t)
	db.Statement.Context = ctx
}

func after(db *gorm.DB) {
	ctx := db.Statement.Context
	stats, ok := cctxv2.GetAs[*CallStats](ctx, cctxv2.MySQLCallStatsKey)
	if !ok {
		return
	}
	var (
		procTime = time.Since(stats.Start)
	)

	var (
		sql    = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars)
		err    = db.Error
		vip, _ = db.Get("vip")
		logMap = map[string]interface{}{
			logger.KParams:   stats.Params,
			logger.KSql:      sql,
			logger.KProcTime: procTime.Milliseconds(),
			logger.KCallee:   vip,
		}
	)
	if err != nil {
		logMap[logger.KErrMsg] = err.Error()
		var (
			isSuccess = errorx.Failed
			errType   = errorx.ErrTypeSys
			code      = errorx.ErrDefault.Code
		)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isSuccess = errorx.Success
			errType = errorx.ErrTypeBiz
			code = errorx.ErrNotFound.Code
		}
		err = errorx.New(errorx.ErrOptions{
			ErrMeta: errorx.ErrMeta{
				ServiceType: errorx.ServiceTypeBasic,
				Service:     errorx.ServiceMysql,
				ErrType:     errType,
				IsSuccess:   isSuccess,
			},
			Code: code,
			Err:  err,
		})
		db.Error = err
	}

	mErr := errorx.Get(err, false)
	if mErr != nil && mErr.ErrType == errorx.ErrTypeSys {
		logger.Warn(ctx, logger.TagMysqlFailure, logMap)
	} else {
		logger.Info(ctx, logger.TagMysqlSuccess, logMap)
	}

}
