package cmd

import (
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/session"
	"github.com/go-macaron/csrf"
	"scheduler/utils"
	"fmt"
	"scheduler/utils/app"
)

// web服务器默认端口
const DefaultPort = 5920
// 静态文件目录
const StaticDir = "public"

var CmdWeb = cli.Command{
	Name: "server",
	Usage: "start scheduler web server",
	Action: run,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name: "port,p",
			Value: DefaultPort,
			Usage: "bind port number",
		},
	},
}

func run(ctx *cli.Context) {
	// 检测环境
	utils.CheckEnv()
	// 启动定时任务
	runScheduler()
	m := macaron.Classic()
	// 注册路由
	registerRouter(m)
	// 注册中间件
	registerMiddleware(m)
	port := parsePort(ctx)
	m.Run(port)
}

// 定时任务调度
func runScheduler()  {}

// 路由注册
func registerRouter(m *macaron.Macaron) {
	// 所有GET方法，自动注册HEAD方法
	m.SetAutoHead(true)
	m.Get("/", func(ctx *macaron.Context) (string) {
		return "go home"
	})
}

// 中间件注册
func registerMiddleware(m *macaron.Macaron)  {
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(gzip.Gziper())
	m.Use(macaron.Static(StaticDir))
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
}

// 解析端口
func parsePort(ctx *cli.Context) int {
	var port int
	if (ctx.IsSet("port")) {
		port = ctx.Int("port")
	}
	if port <= 0 || port >= 65535 {
		port = DefaultPort
	}

	return port
}