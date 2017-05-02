package models

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/core"
    "github.com/go-xorm/xorm"
    "gopkg.in/macaron.v1"
    "strings"
    "time"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/setting"
    "github.com/ouqiang/gocron/modules/app"
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

// 定时ping, 防止因数据库超时设置连接被断开
func keepDbAlived(engine *xorm.Engine)  {
    t := time.Tick(180 * time.Second)
    for {
        <- t
        engine.Ping()
    }
}

// 获取数据库配置
func getDbConfig() map[string]string {
    config, err := setting.Read(app.AppConfig)
    if err != nil {
        logger.Fatal("获取应用配置失败", err)
    }
    section := config.Section("db")
    if err != nil {
        logger.Fatal("获取DB配置失败", err)
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