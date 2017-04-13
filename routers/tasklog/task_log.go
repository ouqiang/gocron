package tasklog

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/utils"
)

// @author qiang.ou<qingqianludao@gmail.com>
// @date 2017/4/7-21:18

func Index(ctx *macaron.Context)  {
    logModel := new(models.TaskLog)
    logs, err := logModel.List()
    if err != nil {
        logger.Error(err)
    }
    ctx.Data["Title"] = "任务日志"
    ctx.Data["Logs"] = logs
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