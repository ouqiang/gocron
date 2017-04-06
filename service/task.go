package service

import (
    "github.com/ouqiang/cron-scheduler/models"
    "github.com/ouqiang/cron-scheduler/modules/ansible"
    "github.com/ouqiang/cron-scheduler/modules/crontask"
    "github.com/robfig/cron"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
    "github.com/ouqiang/cron-scheduler/modules/logger"
)

type Task struct{}

// 初始化任务, 从数据库取出所有任务, 添加到定时任务并运行
func (task *Task) Initialize() {
    taskModel := new(models.Task)
    taskList, err := taskModel.ActiveList()
    if err != nil {
        logger.Error("获取任务列表错误-", err.Error())
        return
    }
    if len(taskList) == 0 {
        logger.Debug("任务列表为空")
        return
    }
    for _, item := range taskList {
        task.Add(item)
    }
}

// 添加任务
func (task *Task) Add(taskModel models.Task) {
    taskFunc := createHandlerJob(taskModel)
    if taskFunc == nil {
        logger.Error("添加任务#不存在的任务协议编号", taskModel.Protocol)
        return
    }
    // 定时任务
    if taskModel.Type == models.Timing {
        err := crontask.DefaultCronTask.Update(strconv.Itoa(taskModel.Id), taskModel.Spec, taskFunc)
        if err != nil {
            logger.Error(err)
        }
    } else if taskModel.Type == models.Delay {
        // 延时任务
        delay := time.Duration(taskModel.Delay) * time.Second
        time.AfterFunc(delay, taskFunc)
    }
}

type Handler interface {
    Run(taskModel models.Task) (string, error)
}

// HTTP任务
type HTTPHandler struct{}

func (h *HTTPHandler) Run(taskModel models.Task) (result string, err error) {
    client := &http.Client{}
    if taskModel.Timeout > 0 {
        client.Timeout = time.Duration(taskModel.Timeout) * time.Second
    }
    req, err := http.NewRequest("POST", taskModel.Command, nil)
    if err != nil {
        logger.Error("创建HTTP请求错误-", err.Error())
        return
    }
    req.Header.Set("Content-type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", "golang/cron-scheduler")

    resp, err := client.Do(req)
    defer func() {
        if resp != nil {
            resp.Body.Close()
        }
    }()
    if err != nil {
        logger.Error("HTTP请求错误-", err.Error())
        return
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Error("读取HTTP请求返回值失败-", err.Error())
    }

    return string(body), err
}

// SSH-command任务
type SSHCommandHandler struct{}

func (ssh *SSHCommandHandler) Run(taskModel models.Task) (string, error) {
    return execSSHHandler("shell", taskModel)
}

// SSH-script任务
type SSHScriptHandler struct{}

func (ssh *SSHScriptHandler) Run(taskModel models.Task) (string, error) {
    return execSSHHandler("script", taskModel)
}

// SSH任务
func execSSHHandler(module string, taskModel models.Task) (string, error) {
    var args []string = []string{taskModel.Command}
    if taskModel.Timeout > 0 {
        // -B 异步执行超时时间, -P 轮询时间
        args = append(args, "-B", strconv.Itoa(taskModel.Timeout), "-P", "10")
    }
    if module == "shell" {
        return ansible.Shell(taskModel.SshHosts, ansible.DefaultHosts.GetFilename(), args...)
    }
    if module == "script" {
        return ansible.Script(taskModel.SshHosts, ansible.DefaultHosts.GetFilename(), args...)
    }

    return "", nil
}

func createTaskLog(taskModel models.Task) (int, error) {
    taskLogModel := new(models.TaskLog)
    taskLogModel.Name = taskModel.Name
    taskLogModel.Spec = taskModel.Spec
    taskLogModel.Protocol = taskModel.Protocol
    taskLogModel.Type = taskModel.Type
    taskLogModel.Command = taskModel.Command
    taskLogModel.Timeout = taskModel.Timeout
    taskLogModel.Delay = taskModel.Delay
    taskLogModel.SshHosts = taskModel.SshHosts
    taskLogModel.StartTime = time.Now()
    taskLogModel.Status = models.Running
    insertId, err := taskLogModel.Create()

    return insertId, err
}

func updateTaskLog(taskLogId int, result string, err error) (int64, error) {
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
    })

}

func createHandlerJob(taskModel models.Task) cron.FuncJob {
    var handler Handler = nil
    switch taskModel.Protocol {
    case models.HTTP:
        handler = new(HTTPHandler)
    case models.SSHCommand:
        handler = new(SSHCommandHandler)
    case models.SSHScript:
        handler = new(SSHScriptHandler)
    }
    if handler == nil {
        return nil
    }
    taskFunc := func() {
        taskLogId, err := createTaskLog(taskModel)
        if err != nil {
            logger.Error("写入任务日志失败-", err)
            return
        }
        // err != nil 执行失败
        result, err := handler.Run(taskModel)
        _, err = updateTaskLog(int(taskLogId), result, err)
        if err != nil {
            logger.Error("更新任务日志失败-", err)
        }
    }

    return taskFunc
}
