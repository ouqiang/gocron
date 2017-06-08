package models

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/core"
    "github.com/go-xorm/xorm"
    "gopkg.in/macaron.v1"
    "strings"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/app"
    "strconv"
    "time"
)

type Status int8
type CommonMap map[string]interface{}

var TablePrefix string = ""
var Db *xorm.Engine

const (
    Disabled Status = 0 // 禁用
    Failure  Status = 0 // 失败
    Enabled  Status = 1 // 启用
    Running  Status = 1 // 运行中
    Finish   Status = 2 // 完成
    Cancel   Status = 3 // 取消
    Waiting Status  = 5 // 等待中
)

const (
    Page        = 1      // 当前页数
    PageSize    = 20     // 每页多少条数据
    MaxPageSize = 100000 // 每次最多取多少条
)

const DefaultTimeFormat = "2006-01-02 15:04:05"

type BaseModel struct  {
    Page      int       `xorm:"-"`
    PageSize  int       `xorm:"-"`
}

func (model *BaseModel) parsePageAndPageSize(params CommonMap) {
    page, ok := params["Page"]
    if ok {
        model.Page = page.(int)
    }
    pageSize, ok := params["PageSize"]
    if ok {
        model.PageSize = pageSize.(int)
    }
    if model.Page <= 0 {
        model.Page = Page
    }
    if model.PageSize <= 0 || model.PageSize > MaxPageSize {
        model.PageSize = PageSize
    }
}

func (model *BaseModel) pageLimitOffset() int {
    return (model.Page - 1) * model.PageSize
}

// 创建Db
func CreateDb() *xorm.Engine {
    config := getDbConfig()
    dsn := getDbEngineDSN(config["engine"], config)
    engine, err := xorm.NewEngine(config["engine"], dsn)
    if err != nil {
        logger.Fatal("创建xorm引擎失败", err)
    }
    maxIdleConns, err := strconv.Atoi(config["max_idle_conns"])
    maxOpenConns, err := strconv.Atoi(config["max_open_conns"])
    if maxIdleConns <= 0 {
        maxIdleConns = 30
    }
    if maxOpenConns <= 0 {
        maxOpenConns = 100
    }
    engine.SetMaxIdleConns(maxIdleConns)
    engine.SetMaxOpenConns(maxOpenConns)

    if config["prefix"] != "" {
        // 设置表前缀
        TablePrefix = config["prefix"]
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


// 获取数据库配置
func getDbConfig() map[string]string {
    var db map[string]string = make(map[string]string)
    db["user"] = app.Setting.Key("db.user").String()
    db["password"] = app.Setting.Key("db.password").String()
    db["host"] = app.Setting.Key("db.host").String()
    db["port"] = app.Setting.Key("db.port").String()
    db["database"] = app.Setting.Key("db.database").String()
    db["charset"] = app.Setting.Key("db.charset").String()
    db["prefix"] = app.Setting.Key("db.prefix").String()
    db["engine"] = app.Setting.Key("db.engine").String()
    db["max_idle_conns"] = app.Setting.Key("db.max.idle.conns").String()
    db["max_open_conns"] = app.Setting.Key("db.max.open.conns").String()

    return db
}

func keepDbAlived(engine *xorm.Engine)  {
    t := time.Tick(180 * time.Second)
    for {
        <- t
        engine.Ping()
    }
}