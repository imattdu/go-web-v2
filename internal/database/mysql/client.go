package mysql

import (
	"context"
	"time"

	"github.com/imattdu/go-web-v2/internal/config"
	"github.com/imattdu/go-web-v2/internal/util/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	GlobalDb *gorm.DB
)

func Init(c context.Context) error {
	if err := NewDb(c); err != nil {
		logger.Error(c, logger.TagUndef, map[string]interface{}{
			"msg": "mysql init failed",
			"err": err.Error(),
		})
		return err
	}
	return nil
}

func NewDb(c context.Context) error {
	dsn := config.GlobalConf.Mysql.StuGoDsn
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn, // DSN data source name
		// DefaultStringSize: 256, // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		//logs.Error(c, logs.LTagUndef).Err(err).Msg("NewStuGoDBCli failed")
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		//logs.Error(c, logs.LTagUndef).Err(err).Msg("stuGoDBCli DB failed")
		return err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	// 空闲连接最大连接数
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// 数据库打开最大连接数
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// 设置连接可以重复使用的最长时间。
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	GlobalDb = db
	_ = GlobalDb.Use(&LogPlugin{})
	return nil
}
