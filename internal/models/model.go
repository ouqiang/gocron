package models

import (
	"fmt"
	"strings"
	"time"

	macaron "gopkg.in/macaron.v1"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/setting"
)

type Status int8
type CommonMap map[string]interface{}

var TablePrefix = ""
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
	Page        = 1    // 当前页数
	PageSize    = 20   // 每页多少条数据
	MaxPageSize = 1000 // 每次最多取多少条
)

const DefaultTimeFormat = "2006-01-02 15:04:05"

const (
	dbPingInterval = 90 * time.Second
	dbMaxLiftTime  = 2 * time.Hour
)

type BaseModel struct {
	Page     int `xorm:"-"`
	PageSize int `xorm:"-"`
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
	if model.PageSize <= 0 {
		model.PageSize = MaxPageSize
	}
}

func (model *BaseModel) pageLimitOffset() int {
	return (model.Page - 1) * model.PageSize
}

// 创建Db
func CreateDb() *xorm.Engine {
	dsn := getDbEngineDSN(app.Setting)
	engine, err := xorm.NewEngine(app.Setting.Db.Engine, dsn)
	if err != nil {
		logger.Fatal("创建xorm引擎失败", err)
	}
	engine.SetMaxIdleConns(app.Setting.Db.MaxIdleConns)
	engine.SetMaxOpenConns(app.Setting.Db.MaxOpenConns)
	engine.SetConnMaxLifetime(dbMaxLiftTime)

	if app.Setting.Db.Prefix != "" {
		// 设置表前缀
		TablePrefix = app.Setting.Db.Prefix
		mapper := core.NewPrefixMapper(core.SnakeMapper{}, app.Setting.Db.Prefix)
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
func CreateTmpDb(setting *setting.Setting) (*xorm.Engine, error) {
	dsn := getDbEngineDSN(setting)

	return xorm.NewEngine(setting.Db.Engine, dsn)
}

// 获取数据库引擎DSN  mysql,sqlite,postgres
func getDbEngineDSN(setting *setting.Setting) string {
	engine := strings.ToLower(setting.Db.Engine)
	dsn := ""
	switch engine {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&allowNativePasswords=true",
			setting.Db.User,
			setting.Db.Password,
			setting.Db.Host,
			setting.Db.Port,
			setting.Db.Database,
			setting.Db.Charset)
	case "postgres":
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			setting.Db.User,
			setting.Db.Password,
			setting.Db.Host,
			setting.Db.Port,
			setting.Db.Database)
	}

	return dsn
}

func keepDbAlived(engine *xorm.Engine) {
	t := time.Tick(dbPingInterval)
	var err error
	for {
		<-t
		err = engine.Ping()
		if err != nil {
			logger.Infof("database ping: %s", err)
		}
	}
}
