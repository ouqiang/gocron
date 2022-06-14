package process

import (
	"fmt"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/gocron/internal/routers/base"
	"github.com/ouqiang/gocron/internal/service"
	"gopkg.in/macaron.v1"
	"strconv"
	"strings"
)

type Form struct {
	Id      int
	Name    string
	Command string
	Tag     string
	NumProc int
	LogFile string
	HostIds string
}

func Index(ctx *macaron.Context) string {
	processModel := new(models.Process)
	queryParams := parseQueryParams(ctx)
	total, err := processModel.Total(queryParams)
	if err != nil {
		logger.Error(err)
	}
	processList, err := processModel.List(queryParams)
	if err != nil {
		logger.Error(err)
	}

	//组装workers数据
	var result []models.Process
	pw := models.ProcessWorker{}
	for _, process := range processList {
		workers, _ := pw.GetByProcess(process)
		process.Workers = workers
		result = append(result, process)
	}

	return utils.JsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  result,
	})
}

func Store(ctx *macaron.Context, form Form) string {
	processModel := models.Process{}

	processModel.Command = form.Command
	processModel.Name = form.Name
	processModel.Tag = form.Tag
	processModel.NumProc = form.NumProc
	processModel.LogFile = form.LogFile

	var err error
	if form.Id == 0 {
		_, err = processModel.Create()
	} else {
		processModel.Id = form.Id
		_, err = processModel.UpdateBean(processModel.Id)
	}
	if err != nil {
		return utils.JsonResp.Failure(500, err.Error())
	}

	pw := models.ProcessWorker{}
	workers, _ := pw.GetByProcess(processModel)
	workerCount := len(workers)
	if workerCount < form.NumProc { //新增worker
		newCount := form.NumProc - workerCount
		for i := 0; i < newCount; i++ {
			pw = models.ProcessWorker{ProcessId: processModel.Id, IsValid: 1}
			err = pw.Create()
			fmt.Println(err)
		}
	}
	fmt.Println(workerCount, workers)
	if workerCount > form.NumProc { //删除worker
		removeCount := workerCount - form.NumProc
		workers, _ := pw.GetLimitByProcess(processModel, removeCount)
		for _, worker := range workers {
			go service.ProcessServiceImpl.StopWorker(worker)
			worker.IsValid = 0
			err = worker.Update()
			logger.Debug(err)
		}
	}

	//处理进程数据
	ph := models.ProcessHost{}
	ph.DeleteForProcess(processModel)
	for _, idStr := range strings.Split(form.HostIds, ",") {
		hostId, _ := strconv.Atoi(idStr)
		if hostId > 0 {
			ph = models.ProcessHost{HostId: hostId, ProcessId: processModel.Id}
			ph.Create()
		}
	}
	return utils.JsonResp.Success("操作成功", processModel)
}

func Get(ctx *macaron.Context) string {
	json := utils.JsonResponse{}
	processId := ctx.ParamsInt("id")
	process := models.Process{}
	err := process.Get(processId)
	if err != nil {
		return json.Failure(400, err.Error())
	}
	ph := models.ProcessHost{}
	process.Hosts = ph.GetByProcess(process)
	return json.Success("Success", process)
}

func Start(ctx *macaron.Context) string {
	id := ctx.ParamsInt("id")
	process := models.Process{}
	_ = process.Get(id)

	if !process.Enable {
		return utils.JsonResp.CommonFailure("进程不可用，请先启动进程", nil)
	}

	// 标记为已开始，后台定时器自动检测并启动关联的worker
	_, err := process.Update(id, models.CommonMap{"status": models.ProcessStart})
	if err != nil {
		return utils.JsonResp.Failure(400, err.Error())
	}
	return utils.JsonResp.Success("Success", nil)
}

func Enable(ctx *macaron.Context) string {
	json := utils.JsonResponse{}
	processId := ctx.ParamsInt("id")
	process := models.Process{}
	err := process.Get(processId)
	if err != nil {
		return json.CommonFailure("获取进程失败", err)
	}
	_, err = process.Update(processId, models.CommonMap{"enable": true})
	if err != nil {
		return json.CommonFailure("更新数据失败", err)
	}
	return json.Success("操作成功", process)
}

func Disable(ctx *macaron.Context) string {
	json := utils.JsonResponse{}
	processId := ctx.ParamsInt("id")
	process := models.Process{}
	err := process.Get(processId)
	if err != nil {
		return json.CommonFailure("获取进程失败", err)
	}
	total, _ := process.GetStartWorkerTotal()
	if total > 0 {
		return json.CommonFailure("进程存在运行中的工作进程,不能禁用,请先停止进程", err)
	}

	_, err = process.Update(processId, models.CommonMap{"enable": false})
	if err != nil {
		return json.CommonFailure("更新数据失败", err)
	}
	return json.Success("操作成功", process)
}

func Stop(ctx *macaron.Context) string {
	json := utils.JsonResponse{}
	id := ctx.ParamsInt("id")
	process := models.Process{}
	_ = process.Get(id)
	process.Update(id, models.CommonMap{"status": models.ProcessStop})
	err := service.ProcessServiceImpl.StopProcess(process)
	if err != nil {
		return json.Failure(400, err.Error())
	}
	return json.Success("Success", nil)
}

func Restart(ctx *macaron.Context) string {
	json := utils.JsonResponse{}
	id := ctx.ParamsInt("id")
	processModel := models.Process{}
	_ = processModel.Get(id)

	return json.Success("Success", nil)
}

func parseQueryParams(ctx *macaron.Context) models.CommonMap {
	var params = models.CommonMap{}
	params["Id"] = ctx.QueryInt("id")
	params["Name"] = ctx.QueryTrim("name")
	params["Command"] = ctx.QueryTrim("command")
	status := ctx.QueryInt("status")
	if status >= 0 {
		status -= 1
	}
	params["Status"] = status
	base.ParsePageAndPageSize(ctx, params)

	return params
}
