package task

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/utils"
    "strings"
)

func Index(ctx *macaron.Context)  {
    taskModel := new(models.Task)
    tasks, err := taskModel.List()
    if err != nil {
        logger.Error(err)
    }
    ctx.Data["Title"] = "任务列表"
    ctx.Data["Tasks"] = tasks
    ctx.Data["URI"] = "/task"
    ctx.HTML(200, "task/index")
}

func Create(ctx *macaron.Context)  {
    hostModel := new(models.Host)
    hosts, err := hostModel.List()
    if err != nil || len(hosts) == 0 {
        logger.Error(err)
        ctx.Redirect("/host/create")
    }
    logger.Debug(hosts)
    ctx.Data["Title"] = "任务管理"
    ctx.Data["Hosts"] = hosts
    ctx.Data["URI"] = "/task/create"
    ctx.Data["FirstHostName"] = hosts[0].Name
    ctx.Data["FirstHostId"] = hosts[0].Id
    ctx.HTML(200, "task/create")
}

type TaskForm struct {
    Name string `binding:"Required"`
    Spec string `binding:"Required"`
    Protocol models.Protocol `binding:"Required"`
    Type models.TaskType `binding:"Required"`
    Command string `binding:"Required"`
    Timeout int
    Delay int
    HostId int16 `binding:"Required;"`
    Remark string
}

// 保存任务
func Store(ctx *macaron.Context, form TaskForm) string  {
    json := utils.Json{}
    taskModel := models.Task{}
    taskModel.Name = form.Name
    taskModel.Spec = strings.Replace(form.Spec, "\n", "|||", 100)
    taskModel.Protocol = form.Protocol
    taskModel.Type = form.Type
    taskModel.Command = form.Command
    taskModel.Timeout = form.Timeout
    taskModel.HostId = form.HostId
    taskModel.Delay = form.Delay
    taskModel.Remark = form.Remark
    _, err := taskModel.Create()
    if err != nil {
        logger.Error(err)
        return json.Failure(utils.ResponseFailure, "保存失败")
    }

    return json.Success("保存成功", nil)
}

// 删除任务
func Remove(ctx *macaron.Context)  {

}