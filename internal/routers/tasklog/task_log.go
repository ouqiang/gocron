package tasklog

// 任务日志

import (
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/gocron/internal/routers/base"
	"github.com/ouqiang/gocron/internal/service"
	"gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) string {
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
	jsonResp := utils.JsonResponse{}

	return jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  logs,
	})
}

// 清空日志
func Clear(ctx *macaron.Context) string {
	taskLogModel := new(models.TaskLog)
	_, err := taskLogModel.Clear()
	json := utils.JsonResponse{}
	if err != nil {
		return json.CommonFailure(utils.FailureContent)
	}

	return json.Success(utils.SuccessContent, nil)
}

// 停止运行中的任务
func Stop(ctx *macaron.Context) string {
	id := ctx.QueryInt64("id")
	taskId := ctx.QueryInt("task_id")
	taskModel := new(models.Task)
	task, err := taskModel.Detail(taskId)
	json := utils.JsonResponse{}
	if err != nil {
		return json.CommonFailure("获取任务信息失败#"+err.Error(), err)
	}
	if task.Protocol != models.TaskRPC {
		return json.CommonFailure("仅支持SHELL任务手动停止")
	}
	if len(task.Hosts) == 0 {
		return json.CommonFailure("任务节点列表为空")
	}
	for _, host := range task.Hosts {
		service.ServiceTask.Stop(host.Name, host.Port, id)

	}

	return json.Success("已执行停止操作, 请等待任务退出", nil)
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

// 解析查询参数
func parseQueryParams(ctx *macaron.Context) models.CommonMap {
	var params models.CommonMap = models.CommonMap{}
	params["TaskId"] = ctx.QueryInt("task_id")
	params["Protocol"] = ctx.QueryInt("protocol")
	status := ctx.QueryInt("status")
	if status >= 0 {
		status -= 1
	}
	params["Status"] = status
	base.ParsePageAndPageSize(ctx, params)

	return params
}
