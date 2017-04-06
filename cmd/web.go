package cmd

import (
    "github.com/go-macaron/csrf"
    "github.com/go-macaron/gzip"
    "github.com/go-macaron/session"
    "github.com/ouqiang/cron-scheduler/modules/app"
    "github.com/ouqiang/cron-scheduler/routers"
    "github.com/urfave/cli"
    "gopkg.in/macaron.v1"
    "os"
    "os/signal"
    "path/filepath"
    "os/exec"
    "syscall"
    "github.com/ouqiang/cron-scheduler/modules/logger"
    "github.com/ouqiang/cron-scheduler/modules/crontask"
)

// 1号进程id
const InitProcess = 1

// web服务器默认端口
const DefaultPort = 5920

// 静态文件目录
const StaticDir = "public"

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
            Value: "dev",
            Usage: "runtime environment, dev|test|prod",
        },
        cli.StringFlag{
            Name: "d",
            Value: "false",
            Usage: "-d=true, run app as daemon",
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
    // 注册中间件
    registerMiddleware(m)
    port := parsePort(ctx)
    m.Run(port)
}

// 中间件注册
func registerMiddleware(m *macaron.Macaron) {
    m.Use(macaron.Logger())
    m.Use(macaron.Recovery())
    m.Use(gzip.Gziper())
    m.Use(macaron.Static(StaticDir))
    m.Use(macaron.Renderer(macaron.RenderOptions{
        Directory:  "templates",
        Extensions: []string{".html"},
        // 模板语法分隔符，默认为 ["{{", "}}"]
        Delims: macaron.Delims{"{{{", "}}}"},
        // 追加的 Content-Type 头信息，默认为 "UTF-8"
        Charset: "UTF-8",
        // 渲染具有缩进格式的 JSON，默认为不缩进
        IndentJSON: true,
        // 渲染具有缩进格式的 XML，默认为不缩进
        IndentXML: true,
    }))
    m.Use(session.Sessioner())
    m.Use(csrf.Csrfer())
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
    var env string = ""
    if ctx.IsSet("env") {
        env = ctx.String("env")
    }

    if env == "prod" {
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
            case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
                // 删除所有任务
                crontask.DefaultCronTask.DeleteAll()
                os.Exit(1)
        }
    }
}

func becomeDaemon(ctx *cli.Context) {
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