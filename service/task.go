package service

import (
    "github.com/ouqiang/gocron/models"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/ssh"
    "github.com/jakecoffman/cron"
    "github.com/ouqiang/gocron/modules/utils"
    "errors"
)

var Cron *cron.Cron

type Task struct{}

// 初始化任务, 从数据库取出所有任务, 添加到定时任务并运行
func (task *Task) Initialize() {
    Cron = cron.New()
    Cron.Start()
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
    task.BatchAdd(taskList)
}

// 批量添加任务
func (task *Task) BatchAdd(tasks []models.TaskHost)  {
    for _, item := range tasks {
        task.Add(item)
    }
}

// 添加任务
func (task *Task) Add(taskModel models.TaskHost) {
    taskFunc := createHandlerJob(taskModel)
    if taskFunc == nil {
        logger.Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
        return
    }

    cronName := strconv.Itoa(taskModel.Id)
    Cron.RemoveJob(cronName)
    err := Cron.AddFunc(taskModel.Spec, taskFunc, cronName)
    if err != nil {
        logger.Error("添加任务到调度器失败#", err)
    }
}

// 直接运行任务
func (task *Task) Run(taskModel models.TaskHost)  {
    go createHandlerJob(taskModel)()
}

type Handler interface {
    Run(taskModel models.TaskHost) (string, error)
}

// 本地命令
type LocalCommandHandler struct {}

func (h *LocalCommandHandler) Run(taskModel models.TaskHost) (string, error)  {
    if taskModel.Command == "" {
        return "", errors.New("invalid command")
    }

    if utils.IsWindows() {
        return h.runOnWindows(taskModel)
    }

    return h.runOnUnix(taskModel)
}

func (h *LocalCommandHandler) runOnWindows(taskModel models.TaskHost) (string, error) {
    outputGBK, err := utils.ExecShellWithTimeout(taskModel.Timeout, "cmd", "/C", taskModel.Command)
    // windows平台编码为gbk，需转换为utf8才能入库
    outputUTF8, ok := utils.GBK2UTF8(outputGBK)
    if ok {
        return outputUTF8, err
    }

    return "命令输出转换编码失败(gbk to utf8)", err
}

func (h *LocalCommandHandler) runOnUnix(taskModel models.TaskHost) (string, error)  {
    return utils.ExecShellWithTimeout(taskModel.Timeout, "/bin/bash", "-c", taskModel.Command)
}

// HTTP任务
type HTTPHandler struct{}

func (h *HTTPHandler) Run(taskModel models.TaskHost) (result string, err error) {
    client := &http.Client{}
    if taskModel.Timeout > 0 {
        client.Timeout = time.Duration(taskModel.Timeout) * time.Second
    }
    req, err := http.NewRequest("GET", taskModel.Command, nil)
    if err != nil {
        logger.Error("任务处理#创建HTTP请求错误-", err.Error())
        return
    }
    req.Header.Set("Content-type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", "golang/gocron")

    resp, err := client.Do(req)
    defer func() {
        if resp != nil {
            resp.Body.Close()
        }
    }()
    if err != nil {
        logger.Error("任务处理HTTP请求错误-", err.Error())
        return
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Error("任务处理#读取HTTP请求返回值失败-", err.Error())
    }

    return string(body), err
}

// SSH-command任务
type SSHCommandHandler struct{}

func (h *SSHCommandHandler) Run(taskModel models.TaskHost) (string, error) {
    sshConfig := ssh.SSHConfig{
        User: taskModel.Username,
        Password: taskModel.Password,
        Host: taskModel.Name,
        Port: taskModel.Port,
        ExecTimeout: taskModel.Timeout,
        AuthType: taskModel.AuthType,
        PrivateKey: taskModel.PrivateKey,
    }
    return ssh.Exec(sshConfig, taskModel.Command)
}

func createTaskLog(taskModel models.TaskHost) (int64, error) {
    taskLogModel := new(models.TaskLog)
    taskLogModel.TaskId = taskModel.Id
    taskLogModel.Name = taskModel.Task.Name
    taskLogModel.Spec = taskModel.Spec
    taskLogModel.Protocol = taskModel.Protocol
    taskLogModel.Command = taskModel.Command
    taskLogModel.Timeout = taskModel.Timeout
    if taskModel.Protocol == models.TaskSSH {
        taskLogModel.Hostname = taskModel.Alias + "-" + taskModel.Name
    }
    taskLogModel.StartTime = time.Now()
    taskLogModel.Status = models.Running
    insertId, err := taskLogModel.Create()

    return insertId, err
}

func updateTaskLog(taskLogId int64, result string, err error) (int64, error) {
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

func createHandlerJob(taskModel models.TaskHost) cron.FuncJob {
    var handler Handler = nil
    switch taskModel.Protocol {
        case models.TaskHTTP:
            handler = new(HTTPHandler)
        case models.TaskSSH:
            handler = new(SSHCommandHandler)
        case models.TaskLocalCommand:
            handler = new(LocalCommandHandler)
    }
    if handler == nil {
        return nil
    }
    taskFunc := func() {
        taskLogId, err := createTaskLog(taskModel)
        if err != nil {
            logger.Error("任务开始执行#写入任务日志失败-", err)
            return
        }
        result, err := handler.Run(taskModel)
        _, err = updateTaskLog(taskLogId, result, err)
        if err != nil {
            logger.Error("任务结束#更新任务日志失败-", err)
        }
    }

    return taskFunc
}
