package install

import (
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/app"
    "github.com/ouqiang/gocron/modules/setting"
    "github.com/ouqiang/gocron/modules/utils"
    "gopkg.in/macaron.v1"
    "strconv"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/go-macaron/binding"
    "fmt"
    "github.com/ouqiang/gocron/service"
)

// 系统安装

type InstallForm struct {
    DbType        string `binding:"In(mysql)"`
    DbHost        string `binding:"Required;MaxSize(50)"`
    DbPort        int    `binding:"Required;Range(1,65535)"`
    DbUsername    string `binding:"Required;MaxSize(50)"`
    DbPassword    string `binding:"Required;MaxSize(30)"`
    DbName        string `binding:"Required;MaxSize(50)"`
    DbTablePrefix string `binding:"MaxSize(20)"`
    AdminUsername string `binding:"Required;MinSize(3)"`
    AdminPassword string `binding:"Required;MinSize(6)"`
    ConfirmAdminPassword string `binding:"Required;MinSize(6)"`
    AdminEmail    string `binding:"Required;Email;MaxSize(50)"`
}

func(f InstallForm) Error(ctx *macaron.Context, errs binding.Errors)  {
    logger.Error(errs)
}

func Create(ctx *macaron.Context) {
    if app.Installed {
        ctx.Redirect("/")
    }
    ctx.Data["Title"] = "安装"
    ctx.Data["DisableNav"] = true
    ctx.HTML(200, "install/create")
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
        return json.CommonFailure("数据库连接失败", err)
    }
    // 写入数据库配置
    err = writeConfig(form)
    if err != nil {
        return json.CommonFailure("数据库配置写入文件失败", err)
    }

    models.Db = models.CreateDb()
    // 创建数据库表
    migration := new(models.Migration)
    err = migration.Exec(form.DbName)
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

    app.Installed = true
    // 初始化定时任务
    serviceTask := new(service.Task)
    serviceTask.Initialize()

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

// 测试数据库连接
func testDbConnection(form InstallForm) error {
    var dbConfig map[string]string = make(map[string]string)
    dbConfig["engine"] = form.DbType
    dbConfig["host"] = form.DbHost
    dbConfig["port"] = strconv.Itoa(form.DbPort)
    dbConfig["user"] = form.DbUsername
    dbConfig["password"] = form.DbPassword
    dbConfig["charset"] = "utf8"
    db, err := models.CreateTmpDb(dbConfig)
    if err != nil {
        return err
    }

    defer  db.Close()
    err = db.Ping()

    return err

}