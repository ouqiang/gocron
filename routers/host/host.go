package host

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/ansible"
)

func Create(ctx *macaron.Context)  {
    ctx.Data["Title"] = "主机管理"
    ctx.HTML(200, "host/create")
}

type HostForm struct {
    Name string `binding:"Required"`
    Alias string `binding:"Required"`
    Username string `binding:"Required"`
    Password string
    Port int `binding:"Required;Range(1-65535)"`
    LoginType int8 `binding:"Required"`
    Remark string `binding:"Required"`
}

func Store(ctx *macaron.Context, form HostForm) string  {
    json := utils.Json{}
    hostModel := new(models.Host)
    hostModel.Name = form.Name
    hostModel.Alias = form.Alias
    hostModel.Username = form.Username
    hostModel.Password = form.Password
    hostModel.Port = form.Port
    hostModel.LoginType = models.LoginType(form.LoginType);
    hostModel.Remark = form.Remark
    _, err := hostModel.Create()
    if err != nil {
        logger.Error(err)
        return json.Failure(utils.ResponseFailure, "保存失败")
    }

    ansible.DefaultHosts.Write()

    return json.Success("保存成功", nil)
}