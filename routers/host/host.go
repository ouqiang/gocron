package host

import "gopkg.in/macaron.v1"

func Create(ctx *macaron.Context)  {
    ctx.Data["Title"] = "主机管理"
    ctx.HTML(200, "host/create")
}

func Store(ctx *macaron.Context)  {

}