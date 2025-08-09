package config

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type Conf struct {
	Age    int     `toml:"age"`
	Server server  `toml:"server"`
	Log    logConf `toml:"log"`
	Mysql  mysql   `toml:"mysql"`
	Kafka  kafka   `toml:"kafka"`
	T      []t     `toml:"t"`
}

type server struct {
	Port string `toml:"port"`
}

type logConf struct {
	Path       string        `toml:"path"`
	FileName   string        `toml:"file_name"`
	MaxBackups int           `toml:"max_backups"`
	Timeout    time.Duration `toml:"timeout"`
}

type mysql struct {
	StuGoDsn       string `toml:"stu_go_dsn"`
	MaxIdleCons    int32  `toml:"max_idle_cons"`
	MaxOpenCons    int32  `toml:"max_open_cons"`
	ConMaxLifeTime int32  `toml:"con_max_lifetime"`
}

type kafka struct {
	Addr  string `toml:"addr"`
	Topic string `toml:"topic"`
}

type t struct {
	A int64 `toml:"a"`
}

var GlobalConf Conf

func Init(path string) error {
	if _, err := toml.DecodeFile(path, &GlobalConf); err != nil {
		log.Println("toml.DecodeFile failed, err=", err.Error())
		return err
	}
	return nil
}
