package install

import (
	"github.com/ouqiang/cron-scheduler/models"
	"github.com/ouqiang/cron-scheduler/modules/app"
	"github.com/ouqiang/cron-scheduler/modules/setting"
	"github.com/ouqiang/cron-scheduler/modules/utils"
	"gopkg.in/macaron.v1"
	"strconv"
)

// 系统安装

type InstallForm struct {
	DbType        string `binding:"IN(mysql)"`
	DbHost        string `binding:"Required"`
	DbPort        int    `binding:"Required"`
	DbUsername    string `binding:"Required"`
	DbPassword    string `binding:"Required"`
	DbName        string `binding:"Required"`
	DbTablePrefix string
	AdminUsername string `binding:"Required;MinSize(3)"`
	AdminPassword string `binding:"Required;MinSize(6)"`
	AdminEmail    string `binding:"Email"`
}

// 显示安装页面
func Show(ctx *macaron.Context) {
	if app.Installed {
		ctx.Redirect("/")
	}
	ctx.Data["Title"] = "安装"
	ctx.Data["DisableNav"] = true
	ctx.HTML(200, "install/show")
}

// 安装,
func Install(ctx *macaron.Context, form InstallForm) string {
	json := utils.Json{}
	if app.Installed {
		return json.Failure(utils.ResponseFailure, "系统已安装成功")
	}
	// 写入数据库配置
	err := writeConfig(form)
	if err != nil {
		utils.RecordLog(err)
		return json.Failure(utils.ResponseFailure, "数据库配置写入文件失败")
	}

	// 初始化Db
	app.InitDb()
	// 创建数据库表
	migration := new(models.Migration)
	err = migration.Exec(form.DbName)
	if err != nil {
		utils.RecordLog(err)
		return json.Failure(utils.ResponseFailure, "创建数据库表失败")
	}

	// 创建管理员账号
	err = createAdminUser(form)
	if err != nil {
		utils.RecordLog(err)
		return json.Failure(utils.ResponseFailure, "创建管理员账号失败")
	}

	// 创建安装锁
	err = app.CreateInstallLock()
	if err != nil {
		utils.RecordLog(err)
		return json.Failure(utils.ResponseFailure, "创建文件安装锁失败")
	}

	// 初始化定时任务等
	app.InitResource()

	return json.Success("安装成功", nil)
}

// 数据库配置写入文件
func writeConfig(form InstallForm) error {
	dbConfig := map[string]map[string]string{
		"db": map[string]string{
			"engine":   form.DbType,
			"host":     form.DbHost,
			"port":     strconv.Itoa(form.DbPort),
			"user":     form.DbUsername,
			"password": form.DbPassword,
			"database": form.DbName,
			"prefix":   form.DbTablePrefix,
			"charset":  "utf8",
		},
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
