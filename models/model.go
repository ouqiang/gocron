package models

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"scheduler/utils/setting"
	"github.com/go-xorm/core"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/macaron.v1"
	"scheduler/utils/app"
)

var Db *xorm.Engine

func init()  {
	if app.Installed {
		Db = createDb()
	}
}

type Status int8
type CommonMap map[string]interface{}

const (
	Disabled Status = 0   // 禁用
	Failure Status  = 0   // 失败
	Enabled Status  = 1   // 启用
	Running Status  = 1   // 运行中
	Finish Status   = 2   // 完成
)

const (
	Page = 1            // 当前页数
	PageSize = 20       // 每页多少条数据
	MaxPageSize = 1000  // 每次最多取多少条
)

// 创建Db
func createDb() *xorm.Engine{
	config,err := setting.Read()
	if err != nil {
		panic(err)
	}
	section := config.Section("db")
	if err != nil {
		panic(err)
	}
	user := section.Key("user").String()
	password := section.Key("password").String()
	host := section.Key("host").String()
	port := section.Key("port").String()
	database := section.Key("database").String()
	charset := section.Key("charset").String()
	prefix := section.Key("prefix").String()

	DSN := "%s:%s@tcp(%s:%s)/%s?charset=%s"
	dsn := fmt.Sprintf(DSN, user, password, host, port, database, charset)
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if prefix != "" {
		// 设置表前缀
		mapper := core.NewPrefixMapper(core.SnakeMapper{}, prefix)
		engine.SetTableMapper(mapper)
	}
	// 本地环境开始日志
	if macaron.Env == macaron.DEV {
		engine.ShowSQL(true)
		engine.Logger().SetLevel(core.LOG_DEBUG)
	}

	return engine
}
