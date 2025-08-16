package mysql

import (
	"context"
	"time"

	"github.com/imattdu/go-web-v2/internal/common/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

var GlobalDb *gorm.DB

func Init(c context.Context) (err error) {
	GlobalDb, err = NewDB(c, config.GlobalConf.Mysql)
	return
}

func NewDB(c context.Context, conf config.MysqlConf) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: conf.Dsn, // DSN data source name
		// DefaultStringSize: 256, // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: gLogger.Default.LogMode(gLogger.Silent), // 关闭日志
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		//logs.Error(c, logs.LTagUndef).Err(err).Msg("stuGoDBCli DB failed")
		return nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	// 空闲连接最大连接数
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// 数据库打开最大连接数
	sqlDB.SetMaxOpenConns(1000)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// 设置连接可以重复使用的最长时间。
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	err = db.Use(&Plugin{
		callee: conf.Vip,
	})
	db = db.Set("vip", conf.Vip)
	return db, err
}
