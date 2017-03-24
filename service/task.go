package service

import (
	"github.com/ouqiang/cron-scheduler/models"
	"github.com/ouqiang/cron-scheduler/modules/utils"
	"net/http"
	"io/ioutil"
	"strconv"
	"time"
	"github.com/ouqiang/cron-scheduler/modules/crontask"
	"github.com/robfig/cron"
	"errors"
)

type Task struct {}

// 初始化任务，从数据库取出所有任务添加到定时任务
func(task *Task) Initialize()  {
	taskModel := new(models.Task)
	taskList, err := taskModel.List()
	if err != nil {
		utils.RecordLog("获取任务列表错误-", err.Error())
		return
	}
	if len(taskList) == 0 {
		utils.RecordLog("任务列表为空")
		return
	}
	for _, item := range(taskList) {
		task.Add(item)
	}
}

// 添加任务
func(task *Task) Add(taskModel models.Task) {
	taskFunc := createHandlerJob(taskModel)
	if taskFunc == nil {
		utils.RecordLog("添加任务#不存在的任务协议编号", taskModel.Protocol)
		return
	}
	// 定时任务
	if taskModel.Type == models.Timing {
		crontask.DefaultCronTask.AddOrReplace(strconv.Itoa(taskModel.Id), taskModel.Spec, taskFunc)
	} else if taskModel.Type == models.Delay {
		// 延时任务
		time.AfterFunc(time.Duration(taskModel.Timeout), taskFunc)
	}
}

type Handler interface {
	Run(taskModel models.Task)
}

// HTTP任务
type HTTPHandler struct {}

func(h *HTTPHandler) Run(taskModel models.Task)  {
	client := &http.Client{}
	if (taskModel.Timeout > 0) {
		client.Timeout = time.Duration(taskModel.Timeout) * time.Second
	}
	req, err := http.NewRequest("POST", taskModel.Command, nil)
	if err != nil {
		utils.RecordLog("创建HTTP请求错误-", err.Error())
		return
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "golang-cron/scheduler")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		utils.RecordLog("HTTP请求错误-", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.RecordLog("读取HTTP请求返回值失败-", err.Error())
	}

}

// SSH任务
type SSHHandler struct {}

func(ssh *SSHHandler) Run(taskModel models.Task)  {

}

func createTaskLog(taskModel models.Task) (int64, error) {
	taskLogModel := new(models.TaskLog)
	taskLogModel.TaskId = taskModel.Id
	taskLogModel.StartTime = time.Now()
	taskLogModel.Status = models.Running
	insertId, err := taskLogModel.Create()

	return insertId, err
}

func updateTaskLog(taskModel models.Task, result string)  {
	taskLogModel := new(models.TaskLog)
	taskLogModel.TaskId= taskModel.Id
	taskLogModel.StartTime = time.Now()

}

func createHandlerJob(taskModel models.Task) cron.FuncJob {
	var taskFunc cron.FuncJob = nil;
	switch taskModel.Protocol {
		case models.HTTP:
			taskFunc = func() {
				var handler Handler = new(HTTPHandler)
				createTaskLog(taskModel)
				handler.Run(taskModel)
			}
		case models.SSH:
			taskFunc = func() {
				var handler Handler = new(SSHHandler)
				createTaskLog(taskModel)
				handler.Run(taskModel)
			}
	}

	return taskFunc
}

