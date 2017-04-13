package models

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/core"
    "github.com/go-xorm/xorm"
    "gopkg.in/macaron.v1"
    "strings"
    "time"
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

const DefaultTimeFormat = "2006-01-02 15:04:05"

// 创建Db
func CreateDb(config map[string]string) *xorm.Engine {
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
    // 本地环境开启日志
    if macaron.Env == macaron.DEV {
        engine.ShowSQL(true)
        engine.Logger().SetLevel(core.LOG_DEBUG)
    }

    go keepDbAlived(engine)

    return engine
}

// 创建临时数据库连接
func CreateTmpDb(config map[string]string) (*xorm.Engine, error)  {
    dsn := getDbEngineDSN(config["engine"], config)

    return xorm.NewEngine(config["engine"], dsn)
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

// 定时ping, 防止因数据库超时设置被断开
func keepDbAlived(engine *xorm.Engine)  {
    t := time.Tick(180 * time.Second)
    for {
        <- t
        engine.Ping()
    }
}