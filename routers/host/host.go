package host

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/ouqiang/gocron/modules/logger"
    "strconv"
)

func Index(ctx *macaron.Context)  {
    hostModel := new(models.Host)
    hosts, err := hostModel.List()
    if err != nil {
        logger.Error(err)
    }
    ctx.Data["Title"] = "主机列表"
    ctx.Data["Hosts"] = hosts
    ctx.HTML(200, "host/index")
}

func Create(ctx *macaron.Context)  {
    ctx.Data["Title"] = "添加主机"
    ctx.HTML(200, "host/host_form")
}

func Edit(ctx *macaron.Context)  {
    ctx.Data["Title"] = "编辑主机"
    hostModel := new(models.Host)
    id := ctx.ParamsInt(":id")
    err := hostModel.Find(id)
    if err != nil {
        logger.Errorf("获取主机详情失败#主机id-%d", id)
    }
    ctx.Data["Host"] = hostModel
    ctx.HTML(200, "host/host_form")
}

type HostForm struct {
    Name string `binding:"Required;MaxSize(100)"`
    Alias string `binding:"Required;MaxSize(32)"`
    Username string `binding:"Required;MaxSize(32)"`
    Password string `binding:"Required;MaxSize(64)"`
    Port int `binding:"Required;Range(1-65535)"`
    Remark string
}

func Store(ctx *macaron.Context, form HostForm) string  {
    json := utils.JsonResponse{}
    hostModel := new(models.Host)
    nameExist, err := hostModel.NameExists(form.Name)
    if err != nil {
        return json.CommonFailure("操作失败", err)
    }
    if nameExist {
        return json.CommonFailure("主机名已存在")
    }
    hostModel.Name = form.Name
    hostModel.Alias = form.Alias
    hostModel.Username = form.Username
    hostModel.Password = form.Password
    hostModel.Port = form.Port
    hostModel.Remark = form.Remark
    _, err = hostModel.Create()
    if err != nil {
        return json.CommonFailure("保存失败", err)
    }

    return json.Success("保存成功", nil)
}

func Remove(ctx *macaron.Context) string  {
    id, err := strconv.Atoi(ctx.Params(":id"))
    json := utils.JsonResponse{}
    if err != nil {
        return json.CommonFailure("参数错误", err)
    }
    taskModel := new(models.Task)
    exist,err := taskModel.HostIdExist(int16(id))
    if err != nil {
        return json.CommonFailure("操作失败", err)
    }
    if exist {
        return json.CommonFailure("有任务引用此主机，不能删除")
    }

    hostModel := new(models.Host)
    _, err =hostModel.Delete(id)
    if err != nil {
        return json.CommonFailure("操作失败", err)
    }

    return json.Success("操作成功", nil)
}