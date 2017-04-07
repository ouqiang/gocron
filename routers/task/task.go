package task

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/utils"
    "strings"
    "github.com/ouqiang/gocron/service"
)

func Create(ctx *macaron.Context)  {
    hostModel := new(models.Host)
    hosts, err := hostModel.List()
    if err != nil || len(hosts) == 0 {
        logger.Error(err)
        ctx.Redirect("/host/create")
    }
    logger.Info(hosts)
    ctx.Data["Title"] = "任务管理"
    ctx.Data["Hosts"] = hosts
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
    Remark string
}

func Store(ctx *macaron.Context, form TaskForm) string  {
    hosts := ctx.Req.Form["hosts[]"]
    taskModel := models.Task{}
    taskModel.Name = form.Name
    taskModel.Spec = form.Spec
    taskModel.Protocol = form.Protocol
    taskModel.Type = form.Type
    taskModel.Command = form.Command
    taskModel.Timeout = form.Timeout
    taskModel.Delay = form.Delay
    taskModel.Remark = form.Remark
    taskModel.SshHosts = strings.Join(hosts, ",")
    _, err := taskModel.Create()
    json := utils.Json{}
    if err != nil {
        logger.Error(err)
        return json.Failure(utils.ResponseFailure, "保存失败")
    }

    serviceTask := new(service.Task)
    serviceTask.Add(taskModel)

    return json.Success("保存成功", nil)
}
