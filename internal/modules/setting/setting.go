package setting

import (
	"errors"

	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"gopkg.in/ini.v1"
)

const DefaultSection = "default"

type Setting struct {
	Db struct {
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
	}
	AllowIps      string
	AppName       string
	ApiKey        string
	ApiSecret     string
	ApiSignEnable bool

	EnableTLS bool
	CAFile    string
	CertFile  string
	KeyFile   string

	ConcurrencyQueue int
	AuthSecret       string
}

// 读取配置
func Read(filename string) (*Setting, error) {
	config, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}
	section := config.Section(DefaultSection)

	var s Setting

	s.Db.Engine = section.Key("db.engine").MustString("mysql")
	s.Db.Host = section.Key("db.host").MustString("127.0.0.1")
	s.Db.Port = section.Key("db.port").MustInt(3306)
	s.Db.User = section.Key("db.user").MustString("")
	s.Db.Password = section.Key("db.password").MustString("")
	s.Db.Database = section.Key("db.database").MustString("gocron")
	s.Db.Prefix = section.Key("db.prefix").MustString("")
	s.Db.Charset = section.Key("db.charset").MustString("utf8")
	s.Db.MaxIdleConns = section.Key("db.max.idle.conns").MustInt(30)
	s.Db.MaxOpenConns = section.Key("db.max.open.conns").MustInt(100)

	s.AllowIps = section.Key("allow_ips").MustString("")
	s.AppName = section.Key("app.name").MustString("定时任务管理系统")
	s.ApiKey = section.Key("api.key").MustString("")
	s.ApiSecret = section.Key("api.secret").MustString("")
	s.ApiSignEnable = section.Key("api.sign.enable").MustBool(true)
	s.ConcurrencyQueue = section.Key("concurrency.queue").MustInt(500)
	s.AuthSecret = section.Key("auth_secret").MustString("")
	if s.AuthSecret == "" {
		s.AuthSecret = utils.RandAuthToken()
	}

	s.EnableTLS = section.Key("enable_tls").MustBool(false)
	s.CAFile = section.Key("ca_file").MustString("")
	s.CertFile = section.Key("cert_file").MustString("")
	s.KeyFile = section.Key("key_file").MustString("")

	if s.EnableTLS {
		if !utils.FileExist(s.CAFile) {
			logger.Fatalf("failed to read ca cert file: %s", s.CAFile)
		}

		if !utils.FileExist(s.CertFile) {
			logger.Fatalf("failed to read client cert file: %s", s.CertFile)
		}

		if !utils.FileExist(s.KeyFile) {
			logger.Fatalf("failed to read client key file: %s", s.KeyFile)
		}
	}

	return &s, nil
}

// 写入配置
func Write(config []string, filename string) error {
	if len(config) == 0 {
		return errors.New("参数不能为空")
	}
	if len(config)%2 != 0 {
		return errors.New("参数不匹配")
	}

	file := ini.Empty()

	section, err := file.NewSection(DefaultSection)
	if err != nil {
		return err
	}
	for i := 0; i < len(config); {
		_, err = section.NewKey(config[i], config[i+1])
		if err != nil {
			return err
		}
		i += 2
	}
	err = file.SaveTo(filename)

	return err
}
