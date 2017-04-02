package models

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/core"
    "github.com/go-xorm/xorm"
    "github.com/ouqiang/cron-scheduler/modules/setting"
    "gopkg.in/macaron.v1"
    "strings"
)

type Status int8
type CommonMap map[string]interface{}

var Db *xorm.Engine

const (
    Disabled Status = 0 // 禁用
    Failure  Status = 0 // 失败
    Enabled  Status = 1 // 启用
    Running  Status = 1 // 运行中
    Finish   Status = 2 // 完成
)

const (
    Page        = 1      // 当前页数
    PageSize    = 20     // 每页多少条数据
    MaxPageSize = 100000 // 每次最多取多少条
)

// 创建Db
func CreateDb(configFile string) *xorm.Engine {
    config := getDbConfig(configFile)
    dsn := getDbEngineDSN(config["engine"], config)
    engine, err := xorm.NewEngine(config["engine"], dsn)
    if err != nil {
        panic(err)
    }
    if config["prefix"] != "" {
        // 设置表前缀
        mapper := core.NewPrefixMapper(core.SnakeMapper{}, config["prefix"])
        engine.SetTableMapper(mapper)
    }
    // 本地环境开始日志
    if macaron.Env == macaron.DEV {
        engine.ShowSQL(true)
        engine.Logger().SetLevel(core.LOG_DEBUG)
    }

    return engine
}

// 获取数据库引擎DSN  mysql,sqlite
func getDbEngineDSN(engine string, config map[string]string) string {
    engine = strings.ToLower(engine)
    var dsn string = ""
    switch engine {
    case "mysql":
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
            config["user"],
            config["password"],
            config["host"],
            config["port"],
            config["database"],
            config["charset"])
    }

    return dsn
}

// 获取数据库配置
func getDbConfig(configFile string) map[string]string {
    config, err := setting.Read(configFile)
    if err != nil {
        panic(err)
    }
    section := config.Section("db")
    if err != nil {
        panic(err)
    }
    var db map[string]string = make(map[string]string)
    db["user"] = section.Key("user").String()
    db["password"] = section.Key("password").String()
    db["host"] = section.Key("host").String()
    db["port"] = section.Key("port").String()
    db["database"] = section.Key("database").String()
    db["charset"] = section.Key("charset").String()
    db["prefix"] = section.Key("prefix").String()
    db["engine"] = section.Key("engine").String()

    return db
}
