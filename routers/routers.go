package routers

import (
    "github.com/go-macaron/binding"
    "github.com/ouqiang/gocron/routers/install"
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/routers/task"
    "github.com/ouqiang/gocron/routers/host"
    "github.com/ouqiang/gocron/routers/tasklog"
)

// 路由注册
func Register(m *macaron.Macaron) {
    // 所有GET方法，自动注册HEAD方法
    m.SetAutoHead(true)
    // 404错误
    m.NotFound(func(ctx *macaron.Context) {
        ctx.HTML(404, "error/404")
    })
    // 50x错误
    m.InternalServerError(func(ctx *macaron.Context) {
        ctx.HTML(500, "error/500")
    })
    // 首页
    m.Get("/", func(ctx *macaron.Context) string {
        return "go home"
    })
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
    })

    // 主机
    m.Group("/host", func() {
        m.Get("/create", host.Create)
        m.Post("/store", binding.Bind(host.HostForm{}), host.Store)
        m.Get("", host.Index)
    })

    // API接口
    m.Group("/api/v1", func() {

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