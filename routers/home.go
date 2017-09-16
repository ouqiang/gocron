package routers

import "gopkg.in/macaron.v1"

// 首页
func Home(ctx *macaron.Context) {
	ctx.Redirect("/task")
}
