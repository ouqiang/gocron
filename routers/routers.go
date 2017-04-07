package routers

import (
    "github.com/go-macaron/binding"
    "github.com/ouqiang/cron-scheduler/routers/install"
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/cron-scheduler/routers/task"
    "github.com/ouqiang/cron-scheduler/routers/host"
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
        m.Post("/store", task.Store)
    })

    // 主机
    m.Group("/host", func() {
        m.Get("/create", host.Create)
        m.Get("/store", host.Store)
    })

    // 任务日志
    m.Group("/tasklog/", func() {

    })

    // API接口
    m.Group("/api/v1", func() {

    })
}
