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
	"github.com/ouqiang/cron-scheduler/modules/ansible"
	"fmt"
)

type Task struct {}

// 初始化任务, 从数据库取出所有任务, 添加到定时任务并运行
func(task *Task) Initialize() {
	taskModel := new(models.Task)
	taskList, err := taskModel.ActiveList()
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
	crontask.DefaultCronTask.Run()
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
		err := crontask.DefaultCronTask.AddOrReplace(strconv.Itoa(taskModel.Id), taskModel.Spec, taskFunc)
		if err != nil {
			utils.RecordLog(err)
		}
	} else if taskModel.Type == models.Delay {
		// 延时任务
		time.AfterFunc(time.Duration(taskModel.Delay) * time.Second, taskFunc)
	}
}

type Handler interface {
	Run(taskModel models.Task) (string, error)
}

// HTTP任务
type HTTPHandler struct {}

func(h *HTTPHandler) Run(taskModel models.Task) (result string, err error) {
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
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		utils.RecordLog("HTTP请求错误-", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.RecordLog("读取HTTP请求返回值失败-", err.Error())
	}

	return string(body),err
}

// SSH任务
type SSHHandler struct {}

func(ssh *SSHHandler) Run(taskModel models.Task) (string, error) {

	var args []string = []string{
		"-m", "shell",
		"-a", taskModel.Command,
	}
	if (taskModel.Timeout > 0) {
		// -B 异步执行超时时间, -P 轮询时间
		args = append(args, "-B", strconv.Itoa(taskModel.Timeout), "-P", "10")
	}
	result, err := ansible.ExecCommand(taskModel.SshHosts, ansible.DefaultHosts.GetFilename(), args...)

	return result, err
}

func createTaskLog(taskId int) (int, error) {
	taskLogModel := new(models.TaskLog)
	taskLogModel.TaskId = taskId
	taskLogModel.StartTime = time.Now()
	taskLogModel.Status = models.Running
	insertId, err := taskLogModel.Create()

	return insertId, err
}

func updateTaskLog(taskLogId int, result string, err error) (int64, error) {
	fmt.Println(taskLogId)
	taskLogModel := new(models.TaskLog)
	var status models.Status
	if err != nil {
		result = err.Error() + " " + result
		status = models.Failure
	} else {
		status = models.Finish
	}
	return taskLogModel.Update(taskLogId, models.CommonMap{
		"status": status,
		"result": result,
	});

}

func createHandlerJob(taskModel models.Task) cron.FuncJob {
	var handler Handler = nil
	switch taskModel.Protocol {
		case models.HTTP:
			handler = new(HTTPHandler)
		case models.SSH:
			handler = new(SSHHandler)
	}
	taskFunc := func() {
		taskLogId, err := createTaskLog(taskModel.Id)
		if err != nil {
			utils.RecordLog("写入任务日志失败-", err)
			return
		}
		// err != nil 执行失败
		result, err := handler.Run(taskModel)
		_, err = updateTaskLog(int(taskLogId), result, err)
		if err != nil {
			utils.RecordLog("更新任务日志失败-", err)
		}
	}

	return taskFunc
}