package routers

import "gopkg.in/macaron.v1"

// 首页
func Home(ctx *macaron.Context)  {
    ctx.Data["Title"] = "首页"
    ctx.HTML(200, "home/index")
}
