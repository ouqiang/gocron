package cmd

import (
    "github.com/ouqiang/gocron/modules/app"
    "github.com/ouqiang/gocron/routers"
    "github.com/urfave/cli"
    "gopkg.in/macaron.v1"
    "os"
    "os/signal"
    "path/filepath"
    "os/exec"
    "syscall"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/utils"
)

// 1号进程id
const InitProcess = 1

// web服务器默认端口
const DefaultPort = 5920

var CmdWeb = cli.Command{
    Name:   "server",
    Usage:  "start scheduler web server",
    Action: run,
    Flags: []cli.Flag{
        cli.IntFlag{
            Name:  "port,p",
            Value: DefaultPort,
            Usage: "bind port number",
        },
        cli.StringFlag{
            Name: "env,e",
            Value: "prod",
            Usage: "runtime environment, dev|test|prod",
        },
        cli.StringFlag{
            Name: "d",
            Value: "false",
            Usage: "-d=true, run app as daemon, not support windows",
        },
    },
}

func run(ctx *cli.Context) {
    // 作为守护进程运行
    becomeDaemon(ctx);
    // 设置运行环境
    setEnvironment(ctx)
    // 初始化应用
    app.InitEnv()
    // 捕捉信号,配置热更新等
    go catchSignal()
    m := macaron.Classic()
    // 注册路由
    routers.Register(m)
    // 注册中间件.
    routers.RegisterMiddleware(m)
    port := parsePort(ctx)
    m.Run(port)
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

func setEnvironment(ctx *cli.Context)  {
    var env string = "prod"
    if ctx.IsSet("env") {
        env = ctx.String("env")
    }

    switch env {
        case "prod":
            macaron.Env = macaron.PROD
        case "test":
            macaron.Env = macaron.TEST
        case "dev":
            macaron.Env = macaron.DEV
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
            case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
                os.Exit(1)
        }
    }
}

// 作为守护进程运行
func becomeDaemon(ctx *cli.Context) {
    // 不支持windows
    if utils.IsWindows() {
        return
    }
    var daemond string = "false"
    if ctx.IsSet("d") {
        daemond = ctx.String("d")
    }
    if (daemond != "true") {
        return
    }

    if os.Getppid() == InitProcess {
        // 已是守护进程，不再处理
        return
    }

    filePath, _:= filepath.Abs(os.Args[0])
    cmd := exec.Command(filePath, os.Args[1:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Start()

    // 父进程退出, 子进程由init-1号进程收养
    os.Exit(0)
}