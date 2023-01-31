package models

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/setting"
	"github.com/stretchr/testify/assert"
)

func TestDSN(t *testing.T) {
	cfg := &setting.Setting{
		Db: struct {
			Engine       string
			Host         string
			Port         int
			User         string
			Password     string
			Database     string
			Prefix       string
			Charset      string
			MaxIdleConns int
			MaxOpenConns int
		}{Engine: "sqlite3", Database: "./test.db"},
		AllowIps:         "",
		AppName:          "",
		ApiKey:           "",
		ApiSecret:        "",
		ApiSignEnable:    false,
		EnableTLS:        false,
		CAFile:           "",
		CertFile:         "",
		KeyFile:          "",
		ConcurrencyQueue: 0,
		AuthSecret:       "",
	}
	dsn := getDbEngineDSN(cfg)
	t.Log("dsn: ", dsn)
	assert.NotEmpty(t, dsn)
}

func TestCreateDb(t *testing.T) {
	logger.InitLogger()
	cfg := &setting.Setting{
		Db: struct {
			Engine       string
			Host         string
			Port         int
			User         string
			Password     string
			Database     string
			Prefix       string
			Charset      string
			MaxIdleConns int
			MaxOpenConns int
		}{Engine: "sqlite3", Database: "./test.db"},
		AllowIps:         "",
		AppName:          "",
		ApiKey:           "",
		ApiSecret:        "",
		ApiSignEnable:    false,
		EnableTLS:        false,
		CAFile:           "",
		CertFile:         "",
		KeyFile:          "",
		ConcurrencyQueue: 0,
		AuthSecret:       "",
	}
	app.Setting = cfg

	os.Remove(cfg.Db.Database)
	defer os.Remove(cfg.Db.Database)

	// 创建数据库并初始化
	Db = CreateDb()
	assert.NotNil(t, Db)

	migrate := new(Migration)
	err := migrate.Install("hello")
	assert.Nil(t, err)
}
