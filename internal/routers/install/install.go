package install

import (
	"errors"
	"fmt"
	"strconv"

	macaron "gopkg.in/macaron.v1"

	"github.com/go-macaron/binding"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/setting"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/gocron/internal/service"
)

// 系统安装

type InstallForm struct {
	DbType               string `binding:"In(mysql,postgres)"`
	DbHost               string `binding:"Required;MaxSize(50)"`
	DbPort               int    `binding:"Required;Range(1,65535)"`
	DbUsername           string `binding:"Required;MaxSize(50)"`
	DbPassword           string `binding:"Required;MaxSize(30)"`
	DbName               string `binding:"Required;MaxSize(50)"`
	DbTablePrefix        string `binding:"MaxSize(20)"`
	AdminUsername        string `binding:"Required;MinSize(3)"`
	AdminPassword        string `binding:"Required;MinSize(6)"`
	ConfirmAdminPassword string `binding:"Required;MinSize(6)"`
	AdminEmail           string `binding:"Required;Email;MaxSize(50)"`
}

func (f InstallForm) Error(ctx *macaron.Context, errs binding.Errors) {
	if len(errs) == 0 {
		return
	}
	json := utils.JsonResponse{}
	content := json.CommonFailure("表单验证失败, 请检测输入")
	ctx.Write([]byte(content))
}

// 安装
func Store(ctx *macaron.Context, form InstallForm) string {
	json := utils.JsonResponse{}
	if app.Installed {
		return json.CommonFailure("系统已安装!")
	}
	if form.AdminPassword != form.ConfirmAdminPassword {
		return json.CommonFailure("两次输入密码不匹配")
	}
	err := testDbConnection(form)
	if err != nil {
		return json.CommonFailure(err.Error())
	}
	// 写入数据库配置
	err = writeConfig(form)
	if err != nil {
		return json.CommonFailure("数据库配置写入文件失败", err)
	}

	appConfig, err := setting.Read(app.AppConfig)
	if err != nil {
		return json.CommonFailure("读取应用配置失败", err)
	}
	app.Setting = appConfig

	models.Db = models.CreateDb()
	// 创建数据库表
	migration := new(models.Migration)
	err = migration.Install(form.DbName)
	if err != nil {
		return json.CommonFailure(fmt.Sprintf("创建数据库表失败-%s", err.Error()), err)
	}

	// 创建管理员账号
	err = createAdminUser(form)
	if err != nil {
		return json.CommonFailure("创建管理员账号失败", err)
	}

	// 创建安装锁
	err = app.CreateInstallLock()
	if err != nil {
		return json.CommonFailure("创建文件安装锁失败", err)
	}

	// 更新版本号文件
	app.UpdateVersionFile()

	app.Installed = true
	// 初始化定时任务
	service.ServiceTask.Initialize()

	return json.Success("安装成功", nil)
}

// 配置写入文件
func writeConfig(form InstallForm) error {
	dbConfig := []string{
		"db.engine", form.DbType,
		"db.host", form.DbHost,
		"db.port", strconv.Itoa(form.DbPort),
		"db.user", form.DbUsername,
		"db.password", form.DbPassword,
		"db.database", form.DbName,
		"db.prefix", form.DbTablePrefix,
		"db.charset", "utf8",
		"db.max.idle.conns", "5",
		"db.max.open.conns", "100",
		"allow_ips", "",
		"app.name", "定时任务管理系统", // 应用名称
		"api.key", "",
		"api.secret", "",
		"enable_tls", "false",
		"concurrency.queue", "500",
		"auth_secret", utils.RandAuthToken(),
		"ca_file", "",
		"cert_file", "",
		"key_file", "",
	}

	return setting.Write(dbConfig, app.AppConfig)
}

// 创建管理员账号
func createAdminUser(form InstallForm) error {
	user := new(models.User)
	user.Name = form.AdminUsername
	user.Password = form.AdminPassword
	user.Email = form.AdminEmail
	user.IsAdmin = 1
	_, err := user.Create()

	return err
}

// 测试数据库连接
func testDbConnection(form InstallForm) error {
	var s setting.Setting
	s.Db.Engine = form.DbType
	s.Db.Host = form.DbHost
	s.Db.Port = form.DbPort
	s.Db.User = form.DbUsername
	s.Db.Password = form.DbPassword
	s.Db.Database = form.DbName
	s.Db.Charset = "utf8"
	db, err := models.CreateTmpDb(&s)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Ping()
	if s.Db.Engine == "postgres" && err != nil {
		pgError, ok := err.(*pq.Error)
		if ok && pgError.Code == "3D000" {
			err = errors.New("数据库不存在")
		}
		return err
	}

	if s.Db.Engine == "mysql" && err != nil {
		mysqlError, ok := err.(*mysql.MySQLError)
		if ok && mysqlError.Number == 1049 {
			err = errors.New("数据库不存在")
		}
		return err
	}

	return err

}
