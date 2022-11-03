// Command gocron
//go:generate statik -src=../../web/public -dest=../../internal -f

package main

import (
	"fmt"
	"github.com/ouqiang/gocron/internal/modules/grace/gracehttp"
	macaron "gopkg.in/macaron.v1"
	"net/http"
	"os"

	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/setting"
	"github.com/ouqiang/gocron/internal/routers"
	"github.com/ouqiang/gocron/internal/service"
	"github.com/ouqiang/goutil"
	"github.com/urfave/cli"
)

var (
	AppVersion           = "1.5"
	BuildDate, GitCommit string
)

// web服务器默认端口
const DefaultPort = 5920

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "gocron"
	cliApp.Usage = "gocron service"
	cliApp.Version, _ = goutil.FormatAppVersion(AppVersion, GitCommit, BuildDate)
	cliApp.Commands = getCommands()
	cliApp.Flags = append(cliApp.Flags, []cli.Flag{}...)
	err := cliApp.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}

// getCommands
func getCommands() []cli.Command {
	command := cli.Command{
		Name:   "web",
		Usage:  "run web server",
		Action: runWeb,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "host",
				Value: "0.0.0.0",
				Usage: "bind host",
			},
			cli.IntFlag{
				Name:  "port,p",
				Value: DefaultPort,
				Usage: "bind port",
			},
			cli.StringFlag{
				Name:  "env,e",
				Value: "prod",
				Usage: "runtime environment, dev|test|prod",
			},
		},
	}

	return []cli.Command{command}
}

func runWeb(ctx *cli.Context) {
	// 设置运行环境
	setEnvironment(ctx)
	// 初始化应用
	app.InitEnv(AppVersion)
	// 初始化模块 DB、定时任务等
	initModule()
	// 捕捉信号,交由grace接管
	// go catchSignal()
	m := macaron.Classic()
	// 注册路由
	routers.Register(m)
	// 注册中间件.
	routers.RegisterMiddleware(m)

	addr := fmt.Sprintf(":%d", parsePort(ctx))
	err := gracehttp.Serve(&http.Server{Addr: addr, Handler: m})
	if err != nil {
		logger.Fatal(err)
	}
}

func initModule() {
	if !app.Installed {
		return
	}

	config, err := setting.Read(app.AppConfig)
	if err != nil {
		logger.Fatal("读取应用配置失败", err)
	}
	app.Setting = config

	// 初始化DB
	models.Db = models.CreateDb()

	// 版本升级
	upgradeIfNeed()

	// 初始化定时任务；交由grace 内部处理
	// service.ServiceTask.Initialize()
}

// 解析端口
func parsePort(ctx *cli.Context) int {
	port := DefaultPort
	if ctx.IsSet("port") {
		port = ctx.Int("port")
	}
	if port <= 0 || port >= 65535 {
		port = DefaultPort
	}

	return port
}

func setEnvironment(ctx *cli.Context) {
	env := "prod"
	if ctx.IsSet("env") {
		env = ctx.String("env")
	}

	switch env {
	case "test":
		macaron.Env = macaron.TEST
	case "dev":
		macaron.Env = macaron.DEV
	default:
		macaron.Env = macaron.PROD
	}
}

// 应用退出
func shutdown() {
	defer func() {
		logger.Info("已退出")
		os.Exit(0)
	}()

	if !app.Installed {
		return
	}
	service.ServiceTask.WaitAndExit()
}

// 判断应用是否需要升级, 当存在版本号文件且版本小于app.VersionId时升级
func upgradeIfNeed() {
	currentVersionId := app.GetCurrentVersionId()
	// 没有版本号文件
	if currentVersionId == 0 {
		return
	}
	if currentVersionId >= app.VersionId {
		return
	}

	migration := new(models.Migration)
	logger.Infof("版本升级开始, 当前版本号%d", currentVersionId)

	migration.Upgrade(currentVersionId)
	app.UpdateVersionFile()

	logger.Infof("已升级到最新版本%d", app.VersionId)
}
