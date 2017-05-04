package tasklog

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/Unknwon/paginater"
    "fmt"
    "html/template"
    "github.com/ouqiang/gocron/routers/base"
)

// @author qiang.ou<qingqianludao@gmail.com>
// @date 2017/4/7-21:18

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