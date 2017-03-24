package routers

import "gopkg.in/macaron.v1"

// 路由注册
func Register(m *macaron.Macaron) {
	// 所有GET方法，自动注册HEAD方法
	m.SetAutoHead(true)
	// 首页
	m.Get("/", func(ctx *macaron.Context) (string) {
		return "go home"
	})
}