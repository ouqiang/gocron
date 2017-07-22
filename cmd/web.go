package cmd

import (
    "github.com/ouqiang/gocron/modules/app"
    "github.com/ouqiang/gocron/routers"
    "github.com/urfave/cli"
    "gopkg.in/macaron.v1"
    "os"
    "os/signal"
    "syscall"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/service"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/setting"
    "time"
    "github.com/ouqiang/gocron/modules/rpc/grpcpool"
)

// web服务器默认端口
const DefaultPort = 5920


var CmdWeb = cli.Command{
    Name:   "web",
    Usage:  "run web server",
    Action: runWeb,
    Flags: []cli.Flag{
        cli.StringFlag{
            Name: "host",
            Value: "0.0.0.0",
            Usage: "bind host",
        },
        cli.IntFlag{
            Name:  "port,p",
            Value: DefaultPort,
            Usage: "bind port",
        },
        cli.StringFlag{
            Name: "env,e",
            Value: "prod",
            Usage: "runtime environment, dev|test|prod",
        },
        cli.BoolFlag{
            Name: "d",
            Usage: "-d=true, run as daemon process",
        },
    },
}

func runWeb(ctx *cli.Context) {
    // 设置运行环境
    setEnvironment(ctx)
    // 初始化应用
    app.InitEnv()
    // 初始化模块 DB、定时任务等
    initModule()
    // 捕捉信号,配置热更新等
    go catchSignal()
    m := macaron.Classic()

    // 注册路由
    routers.Register(m)
    // 注册中间件.
    routers.RegisterMiddleware(m)
    host := parseHost(ctx)
    port := parsePort(ctx)
    m.Run(host, port)
}

func initModule()  {
    if !app.Installed {
        return
    }

    config, err := setting.Read(app.AppConfig)
    if err != nil {
        logger.Fatal("读取应用配置失败", err)
    }
    app.Setting = config

    models.Db = models.CreateDb()

    // 初始化定时任务
    serviceTask := new(service.Task)
    serviceTask.Initialize()
}

// 解析端口
func parsePort(ctx *cli.Context) int {
    var port int = DefaultPort
    if ctx.IsSet("port") {
        port = ctx.Int("port")
    }
    if port <= 0 || port >= 65535 {
        port = DefaultPort
    }

    return port
}

func parseHost(ctx *cli.Context) string  {
    if ctx.IsSet("host") {
        return ctx.String("host")
    }

    return "0.0.0.0"
}

func setEnvironment(ctx *cli.Context)  {
    var env string = "prod"
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

// 捕捉信号
func catchSignal()  {
    c := make(chan os.Signal)
    // todo 配置热更新, windows 不支持 syscall.SIGUSR1, syscall.SIGUSR2
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
    for {
        s := <- c
        logger.Info("收到信号 -- ", s)
        switch s {
            case syscall.SIGHUP:
            logger.Info("收到终端断开信号, 忽略")
            case  syscall.SIGINT, syscall.SIGTERM:
            shutdown()
        }
    }
}


// 应用退出
func shutdown()  {
    defer func() {
        logger.Info("已退出")
        os.Exit(0)
    }()

    if !app.Installed {
        return
    }
    logger.Info("应用准备退出")
    serviceTask := new(service.Task)
    // 停止所有任务调度
    logger.Info("停止定时任务调度")
    serviceTask.StopAll()
    // 释放gRPC连接池
    grpcpool.Pool.ReleaseAll()

    taskNumInRunning := service.TaskNum.Num()
    logger.Infof("正在运行的任务有%d个", taskNumInRunning)
    if taskNumInRunning > 0 {
        logger.Info("等待所有任务执行完成后退出")
    }
    for {
        if taskNumInRunning <= 0 {
            break
        }
        time.Sleep(3 * time.Second)
        taskNumInRunning = service.TaskNum.Num()
    }
}