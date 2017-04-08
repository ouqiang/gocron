package tasklog

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/logger"
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
    ctx.Data["URI"] = "/task/log"
    ctx.HTML(200, "task/log")
}