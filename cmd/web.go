package cmd

import (
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/session"
	"github.com/ouqiang/cron-scheduler/modules/app"
	"github.com/ouqiang/cron-scheduler/routers"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
)

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
	},
}

func run(ctx *cli.Context) {
	app.InitEnv()
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
	var port int
	if ctx.IsSet("port") {
		port = ctx.Int("port")
	}
	if port <= 0 || port >= 65535 {
		port = DefaultPort
	}

	return port
}
