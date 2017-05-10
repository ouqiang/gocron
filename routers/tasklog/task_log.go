package tasklog

// 任务日志

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/Unknwon/paginater"
    "fmt"
    "html/template"
    "github.com/ouqiang/gocron/routers/base"
    "github.com/ouqiang/gocron/service"
    "errors"
)

func Index(ctx *macaron.Context)  {
    logModel := new(models.TaskLog)
    queryParams := parseQueryParams(ctx)
    total, err := logModel.Total(queryParams)
    if err != nil {
        logger.Error(err)
    }
    logs, err := logModel.List(queryParams)
    if err != nil {
        logger.Error(err)
    }
    PageParams := fmt.Sprintf("task_id=%d&protocol=%d&status=%d&page_size=%d",
        queryParams["TaskId"], queryParams["Protocol"], queryParams["Status"],
        queryParams["PageSize"]);
    queryParams["PageParams"] = template.URL(PageParams)
    p := paginater.New(int(total), queryParams["PageSize"].(int), queryParams["Page"].(int), 5)
    ctx.Data["Pagination"] = p
    ctx.Data["Title"] = "任务日志"
    ctx.Data["Logs"] = logs
    ctx.Data["Params"] = queryParams
    ctx.HTML(200, "task/log")
}

// 清空日志
func Clear(ctx *macaron.Context) string  {
    taskLogModel := new(models.TaskLog)
    _, err := taskLogModel.Clear()
    json := utils.JsonResponse{}
    if err != nil {
        return json.CommonFailure(utils.FailureContent)
    }

    return json.Success(utils.SuccessContent, nil)
}

// 删除N个月前的日志
func Remove(ctx *macaron.Context) string {
    month := ctx.ParamsInt(":id")
    json := utils.JsonResponse{}
    if month < 1 || month > 12 {
        return json.CommonFailure("参数取值范围1-12")
    }
    taskLogModel := new(models.TaskLog)
    _, err := taskLogModel.Remove(month)
    if err != nil {
        return json.CommonFailure("删除失败", err)
    }

    return json.Success("删除成功", nil)
}

// 更新任务状态
func UpdateStatus(ctx *macaron.Context) string {
    id := ctx.QueryTrim("id")
    status := ctx.QueryInt("status")
    result := ctx.QueryTrim("result")
    json := utils.JsonResponse{}

    if id == "" {
        return json.CommonFailure("任务ID不能为空")
    }
    if status != 1 && status != 2 {
        return json.CommonFailure("status值错误")
    }
    if status == 1 {
        status -= 1
    }
    taskLogModel := new(models.TaskLog)
    affectRows, err := taskLogModel.UpdateStatus(id, models.Status(status), result)
    if err != nil || affectRows == 0 {
        return json.CommonFailure("更新任务状态失败")
    }

    // 发送通知
    taskId, err := taskLogModel.GetTaskIdByNotifyId(id)
    if err != nil || taskId <= 0 {
        logger.Error("异步任务回调#根据notify-id获取taskId失败", err)
        return json.Success("success", nil)
    }
    taskModel := new(models.Task)
    task, err := taskModel.Detail(taskId)
    if err != nil || task.Id <= 0 {
        logger.Error("异步任务回调#根据获取任务详情失败", err)
        return json.Success("success", nil)
    }

    taskResult := service.TaskResult{}
    taskResult.Result = result
    if status == 0 {
        taskResult.Err = errors.New("error")
    }
    service.SendNotification(task, taskResult)

    return json.Success("success", nil)
}

// 解析查询参数
func parseQueryParams(ctx *macaron.Context) (models.CommonMap) {
    var params models.CommonMap = models.CommonMap{}
    params["TaskId"] = ctx.QueryInt("task_id")
    params["Protocol"] = ctx.QueryInt("protocol")
    status := ctx.QueryInt("status")
    if status >=0 {
        status -= 1
    }
    params["Status"] = status
    base.ParsePageAndPageSize(ctx, params)

    return params
}