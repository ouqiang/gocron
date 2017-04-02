package routers

import (
    "github.com/go-macaron/binding"
    "github.com/ouqiang/cron-scheduler/routers/install"
    "gopkg.in/macaron.v1"
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
        m.Get("", install.Show)
        m.Post("", binding.Bind(install.InstallForm{}), install.Install)
    })
}
