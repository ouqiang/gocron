package task

import "gopkg.in/macaron.v1"

func Create(ctx *macaron.Context)  {
    ctx.Data["Title"] = "任务管理"
    ctx.HTML(200, "task/create")
}

func Store(ctx *macaron.Context)  {

}