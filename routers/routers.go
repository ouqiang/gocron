package routers

import (
    "github.com/go-macaron/binding"
    "github.com/ouqiang/gocron/routers/install"
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/routers/task"
    "github.com/ouqiang/gocron/routers/host"
    "github.com/ouqiang/gocron/routers/tasklog"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/go-macaron/session"
    "github.com/go-macaron/csrf"
    "github.com/go-macaron/toolbox"
    "github.com/go-macaron/gzip"
    "strings"
    "github.com/ouqiang/gocron/modules/app"
)

// 静态文件目录
const StaticDir = "public"

// 路由注册
func Register(m *macaron.Macaron) {
    // 所有GET方法，自动注册HEAD方法
    m.SetAutoHead(true)
    // 404错误
    m.NotFound(func(ctx *macaron.Context) {
        if isGetRequest(ctx) && !isAjaxRequest(ctx) {
            ctx.Data["Title"] = "404 - NOT FOUND"
            ctx.HTML(404, "error/404")
        } else {
            json := utils.JsonResponse{}
            ctx.Resp.Write([]byte(json.Failure(utils.NotFound, "您访问的地址不存在")))
        }
    })
    // 50x错误
    m.InternalServerError(func(ctx *macaron.Context) {
        if isGetRequest(ctx) && !isAjaxRequest(ctx) {
            ctx.Data["Title"] = "500 - SERVER INTERNAL ERROR"
            ctx.HTML(500, "error/500")
        } else {
            json := utils.JsonResponse{}
            ctx.Resp.Write([]byte(json.Failure(utils.ServerError, "网站暂时无法访问,请稍后再试")))
        }
    })
    // 首页
    m.Get("/", Home)
    // 系统安装
    m.Group("/install", func() {
        m.Get("", install.Create)
        m.Post("/store", binding.Bind(install.InstallForm{}), install.Store)
    })

    // 用户
    m.Group("/user", func() {

    })

    // 任务
    m.Group("/task", func() {
        m.Get("/create", task.Create)
        m.Post("/store", binding.Bind(task.TaskForm{}), task.Store)
        m.Get("", task.Index)
        m.Get("/log", tasklog.Index)
        m.Post("/log/clear", tasklog.Clear)
        m.Post("/remove/:id", task.Remove)
        m.Post("/enable/:id", task.Enable)
        m.Post("/disable/:id", task.Disable)
    })

    // 主机
    m.Group("/host", func() {
        m.Get("/create", host.Create)
        m.Post("/store", binding.Bind(host.HostForm{}), host.Store)
        m.Get("", host.Index)
        m.Post("/remove/:id", host.Remove)
    })
}

// 中间件注册
func RegisterMiddleware(m *macaron.Macaron) {
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
    m.Use(toolbox.Toolboxer(m))

    // 系统未安装，重定向到安装页面
    m.Use(func(ctx *macaron.Context) {
        installUrl := "/install"
        if strings.HasPrefix(ctx.Req.RequestURI, installUrl) {
            return
        }
        if !app.Installed {
            ctx.Redirect(installUrl)
        }
    })
    // 设置模板共享变量
    m.Use(func(ctx *macaron.Context) {
        ctx.Data["URI"] = ctx.Req.RequestURI
        ctx.Data["StandardTimeFormat"] = "2006-01-02 15:03:04"
    })
}

func isAjaxRequest(ctx *macaron.Context) bool {
    req := ctx.Req.Header.Get("X-Requested-With")
    if req == "XMLHttpRequest" {
        return true
    }

    return false
}

func isGetRequest(ctx *macaron.Context) bool {
    return ctx.Req.Method == "GET"
}