package delaytask

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/utils"
    "strings"
    "github.com/ouqiang/gocron/service"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/Unknwon/paginater"
    "fmt"
    "github.com/ouqiang/gocron/routers/base"
    "html/template"
    "github.com/ouqiang/gocron/modules/app"
)

func Index(ctx *macaron.Context)  {
    delayTaskModel := new(models.DelayTask)
    queryParams := parseQueryParams(ctx)
    total, err := delayTaskModel.Total(queryParams)
    tasks, err := delayTaskModel.List(queryParams)
    if err != nil {
        logger.Error(err)
    }
    PageParams := fmt.Sprintf("status=%d&page_size=%d",
        queryParams["Status"], queryParams["PageSize"]);
    queryParams["PageParams"] = template.URL(PageParams)
    p := paginater.New(int(total), queryParams["PageSize"].(int), queryParams["Page"].(int), 5)
    ctx.Data["Pagination"] = p
    ctx.Data["Title"] = "延时任务列表"
    ctx.Data["Tasks"] = tasks
    ctx.Data["Params"] = queryParams
    ctx.HTML(200, "task/delay_task")
}

func Create(ctx *macaron.Context) string {
    url := ctx.QueryTrim("url")
    params := ctx.QueryTrim("params")
    delay := ctx.QueryInt("delay")
    json := utils.JsonResponse{}
    delayTaskEnabled, _ := app.Setting.Key("delay.task.enable").Bool()
    if !delayTaskEnabled {
        return json.CommonFailure("未开启延时任务")
    }
    if url == ""  {
        return json.CommonFailure("url地址不能为空")
    }
    lowerUrl := strings.ToLower(url)
    if !strings.HasPrefix(lowerUrl, "http") &&
        !strings.HasPrefix(lowerUrl, "https") {
        return json.CommonFailure("无效的url地址")
    }
    if len(url) > 128 {
        return json.CommonFailure("url长度不能超过128")
    }
    maxDelay := 1 << 31
    if delay <= 0 || delay > maxDelay {
        return json.CommonFailure("无效的delay, 取值范围1-(2^31-1)")
    }
    if len(params) > 256 {
        return json.CommonFailure("params长度不能超过256")
    }

    delayTask := new(models.DelayTask)
    delayTask.Url = url
    delayTask.Params = params
    delayTask.Delay = delay
    delayTask.Status = models.Waiting
    _, err := delayTask.Create()

    if err != nil {
        return json.CommonFailure("添加失败", err)
    }

    logger.Infof("新增延时任务#id-%d#url-%s#params-%s#delay-%d",
        delayTask.Id, delayTask.Url, delayTask.Params, delayTask.Delay)
    delayTaskService := new(service.DelayTask)
    delayTaskService.Add(*delayTask)

    return json.Success("添加成功", nil)
}

// 删除N个月前的日志
func Remove(ctx *macaron.Context) string {
    month := ctx.ParamsInt(":id")
    json := utils.JsonResponse{}
    if month < 1 || month > 12 {
        return json.CommonFailure("参数取值范围1-12")
    }
    delayTaskModel := new(models.DelayTask)
    _, err := delayTaskModel.Remove(month)
    if err != nil {
        return json.CommonFailure("删除失败", err)
    }

    return json.Success("删除成功", nil)
}

// 解析查询参数
func parseQueryParams(ctx *macaron.Context) (models.CommonMap) {
    var params models.CommonMap = models.CommonMap{}
    status := ctx.QueryInt("status")
    if status >=0 {
        status -= 1
    }
    params["Status"] = status
    base.ParsePageAndPageSize(ctx, params)

    return params
}