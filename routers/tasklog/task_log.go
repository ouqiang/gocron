package tasklog

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/Unknwon/paginater"
    "fmt"
    "html/template"
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
    page := ctx.QueryInt("page")
    pageSize := ctx.QueryInt("page_size")
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = models.PageSize
    }

    params["Page"] = page
    params["PageSize"] = pageSize

    return params
}