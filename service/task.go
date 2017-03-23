package service

import (
	"scheduler/models"
	"scheduler/modules/utils"
	"net/http"
	"io/ioutil"
	"strconv"
	"time"
	"scheduler/modules/crontask"
)

type Task struct {}

// 初始化任务，从数据库取出所有任务添加到定时任务
func(task *Task) Initialize()  {
	taskModel := new(models.Task)
	taskList, err := taskModel.List()
	if err != nil {
		utils.RecordLog("获取任务列表错误-", err.Error())
	}
	if len(taskList) == 0 {
		utils.RecordLog("任务列表为空")
	}
	for _, item := range(taskList) {
		task.Add(item)
	}
}

// 添加任务
func(task *Task) Add(taskModel models.Task)  {
	var taskFunc func() = nil;
	switch taskModel.Protocol {
		case models.HTTP:
			taskFunc = func() {
				var handler Handler = new(HTTPHandler)
				handler.Run(taskModel)
			}
		case models.SSH:
			taskFunc = func() {
				var handler Handler = new(SSHHandler)
				handler.Run(taskModel)
			}
		default:
			utils.RecordLog("任务协议不存在-协议编号: ", taskModel.Protocol)
	}
	if (taskFunc != nil) {
		crontask.DefaultCronTask.Add(strconv.Itoa(taskModel.Id), taskModel.Spec, taskFunc)
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

	_, err = taskModel.Update(
		taskModel.Id,
		models.CommonMap{
			"status": 0,
			"result" : string(body),
		});
	if err != nil {
		utils.RecordLog("更新任务日志失败-", err.Error())
	}
}

// SSH任务
type SSHHandler struct {}

func(ssh *SSHHandler) Run(taskModel models.Task)  {

}

// 延时任务
type DelayHandler struct {}

func (handler *DelayHandler) Run(taskModel models.Task)  {

}


