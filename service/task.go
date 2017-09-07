package service

import (
    "github.com/ouqiang/gocron/models"
    "strconv"
    "time"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/jakecoffman/cron"
    "errors"
    "fmt"
    "github.com/ouqiang/gocron/modules/httpclient"
    "github.com/ouqiang/gocron/modules/notify"
    "sync"
    rpcClient "github.com/ouqiang/gocron/modules/rpc/client"
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "strings"
    "text/template"
    "bytes"
    "encoding/base64"
)

// 定时任务调度管理器
var Cron *cron.Cron
// 同一任务是否有实例处于运行中
var runInstance Instance
// 任务计数-正在运行中的任务
var TaskNum TaskCount

// 任务计数
type TaskCount struct {
    num int
    sync.RWMutex
}

func (c *TaskCount) Add()  {
    c.Lock()
    defer c.Unlock()
    c.num += 1
}

func (c *TaskCount) Done()  {
    c.Lock()
    defer c.Unlock()
    c.num -= 1
}

func (c *TaskCount) Num() int  {
    c.RLock()
    defer c.RUnlock()

    return c.num
}

// 任务ID作为Key
type Instance struct {
    Status map[int]bool
    sync.RWMutex
}

// 是否有任务处于运行中
func (i *Instance) has(key int) bool {
    i.RLock()
    defer i.RUnlock()
    running, ok := i.Status[key]
    if ok && running {
        return true
    }

    return false
}

func (i *Instance) add(key int)  {
    i.Lock()
    defer i.Unlock()
    i.Status[key] = true
}

func (i *Instance) done(key int)  {
    i.Lock()
    defer i.Unlock()
    delete(i.Status, key)
}

type Task struct{}

type TaskResult struct {
    Result string
    Err error
    RetryTimes int8
}

// 初始化任务, 从数据库取出所有任务, 添加到定时任务并运行
func (task *Task) Initialize() {
    Cron = cron.New()
    Cron.Start()
    runInstance = Instance{make(map[int]bool), sync.RWMutex{}}
    TaskNum = TaskCount{0, sync.RWMutex{}}

    taskModel := new(models.Task)
    taskList, err := taskModel.ActiveList()
    if err != nil {
        logger.Error("定时任务初始化#获取任务列表错误-", err.Error())
        return
    }
    if len(taskList) == 0 {
        logger.Debug("任务列表为空")
        return
    }
    task.BatchAdd(taskList)
}

// 批量添加任务
func (task *Task) BatchAdd(tasks []models.Task)  {
    for _, item := range tasks {
        task.Add(item)
    }
}

// 添加任务
func (task *Task) Add(taskModel models.Task) {
    if taskModel.Level == models.TaskLevelChild {
        logger.Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", taskModel.Id);
        return
    }
    taskFunc := createJob(taskModel)
    if taskFunc == nil {
        logger.Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
        return
    }

    cronName := strconv.Itoa(taskModel.Id)
    // Cron任务采用数组存储, 删除任务需遍历数组, 并对数组重新赋值, 任务较多时，有性能问题
    Cron.RemoveJob(cronName)
    err := Cron.AddFunc(taskModel.Spec, taskFunc, cronName)
    if err != nil {
        logger.Error("添加任务到调度器失败#", err)
    }
}

// 停止所有任务
func (task *Task) StopAll()  {
    Cron.Stop()
}

// 直接运行任务
func (task *Task) Run(taskModel models.Task)  {
    go createJob(taskModel)()
}

type Handler interface {
    Run(taskModel models.Task) (string, error)
}


// HTTP任务
type HTTPHandler struct{}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(taskModel models.Task) (result string, err error) {
    if taskModel.Timeout <= 0 || taskModel.Timeout > HttpExecTimeout {
        taskModel.Timeout = HttpExecTimeout
    }
    resp := httpclient.Get(taskModel.Command, taskModel.Timeout)
    // 返回状态码非200，均为失败
    if resp.StatusCode != 200 {
        return resp.Body, errors.New(fmt.Sprintf("HTTP状态码非200-->%d", resp.StatusCode))
    }

    return resp.Body, err
}

// RPC调用执行任务
type RPCHandler struct {}

func (h *RPCHandler) Run(taskModel models.Task) (result string, err error)  {
    taskRequest := new(pb.TaskRequest)
    taskRequest.Timeout = int32(taskModel.Timeout)
    taskRequest.Command = taskModel.Command
    var resultChan chan TaskResult = make(chan TaskResult, len(taskModel.Hosts))
    for _, taskHost := range taskModel.Hosts {
        go func(th models.TaskHostDetail) {
            output, err := rpcClient.ExecWithRetry(th.Name, th.Port, taskRequest)
            var errorMessage string = ""
            if err != nil {
                errorMessage = err.Error()
            }
            outputMessage := fmt.Sprintf("主机: [%s-%s]\n%s\n%s\n\n",
                th.Alias, th.Name, errorMessage, output,
            )
            resultChan <- TaskResult{Err:err, Result: outputMessage}
        }(taskHost)
    }

    var aggregationErr error = nil
    var aggregationResult string = ""
    for i := 0; i < len(taskModel.Hosts); i++ {
        taskResult := <- resultChan
        aggregationResult += taskResult.Result
        if taskResult.Err != nil {
            aggregationErr = taskResult.Err
        }
    }

    return aggregationResult, aggregationErr
}


// 创建任务日志
func createTaskLog(taskModel models.Task, status models.Status) (int64, error) {
    taskLogModel := new(models.TaskLog)
    taskLogModel.TaskId = taskModel.Id
    taskLogModel.Name = taskModel.Name
    taskLogModel.Spec = taskModel.Spec
    taskLogModel.Protocol = taskModel.Protocol
    taskLogModel.Command = taskModel.Command
    taskLogModel.Timeout = taskModel.Timeout
    if taskModel.Protocol == models.TaskRPC {
        var aggregationHost string = ""
        for _, host := range taskModel.Hosts {
            aggregationHost += fmt.Sprintf("%s-%s<br>", host.Alias, host.Name)
        }
        taskLogModel.Hostname = aggregationHost
    }
    taskLogModel.StartTime = time.Now()
    taskLogModel.Status = status
    insertId, err := taskLogModel.Create()

    return insertId, err
}

// 更新任务日志
func updateTaskLog(taskLogId int64, taskResult TaskResult) (int64, error) {
    taskLogModel := new(models.TaskLog)
    var status models.Status
    var result string = taskResult.Result
    if taskResult.Err != nil {
        status = models.Failure
    }  else {
        status = models.Finish
    }
    return taskLogModel.Update(taskLogId, models.CommonMap{
        "retry_times": taskResult.RetryTimes,
        "status": status,
        "result": result,
    })

}

func createJob(taskModel models.Task) cron.FuncJob {
    var handler Handler = createHandler(taskModel)
    if handler == nil {
        return nil
    }
    taskFunc := func() {
        TaskNum.Add()
        defer TaskNum.Done()
        taskLogId := beforeExecJob(taskModel)
        if taskLogId <= 0 {
            return
        }
        logger.Infof("开始执行任务#%s#命令-%s", taskModel.Name, taskModel.Command)
        taskResult := execJob(handler, taskModel)
        logger.Infof("任务完成#%s#命令-%s", taskModel.Name, taskModel.Command)
        afterExecJob(taskModel, taskResult, taskLogId)
    }

    return taskFunc
}

func createHandler(taskModel models.Task) Handler  {
    var handler Handler = nil
    switch taskModel.Protocol {
        case models.TaskHTTP:
            handler = new(HTTPHandler)
        case models.TaskRPC:
            handler = new(RPCHandler)
    }


    return handler;
}

// 任务前置操作
func beforeExecJob(taskModel models.Task) (taskLogId int64)  {
    if taskModel.Multi == 0 && runInstance.has(taskModel.Id) {
        createTaskLog(taskModel, models.Cancel)
        return
    }
    if taskModel.Multi == 0 {
        runInstance.add(taskModel.Id)
    }
    taskLogId, err := createTaskLog(taskModel, models.Running)
    if err != nil {
        logger.Error("任务开始执行#写入任务日志失败-", err)
        return
    }

    logger.Debugf("任务命令-%s", taskModel.Command)

    return taskLogId
}

// 任务执行后置操作
func afterExecJob(taskModel models.Task, taskResult TaskResult, taskLogId int64)  {
    if taskResult.Err != nil {
        taskResult.Result = taskResult.Err.Error() + "\n" + taskResult.Result
    }
    _, err := updateTaskLog(taskLogId, taskResult)
    if err != nil {
        logger.Error("任务结束#更新任务日志失败-", err)
    }

    // 发送邮件
    go SendNotification(taskModel, taskResult)
    // 执行依赖任务
    go execDependencyTask(taskModel, taskResult)
}

// 执行依赖任务, 多个任务并发执行
func execDependencyTask(taskModel models.Task, taskResult TaskResult)  {
    // 父任务才能执行子任务
    if taskModel.Level != models.TaskLevelParent {
        return
    }

    // 是否存在子任务
    dependencyTaskId := strings.TrimSpace(taskModel.DependencyTaskId)
    if dependencyTaskId == "" {
        return
    }

    // 父子任务关系为强依赖, 父任务执行失败, 不执行依赖任务
    if taskModel.DependencyStatus == models.TaskDependencyStatusStrong && taskResult.Err != nil {
        logger.Infof("父子任务为强依赖关系, 父任务执行失败, 不运行依赖任务#主任务ID-%d", taskModel.Id)
        return
    }

    // 获取子任务
    model := new(models.Task)
    tasks , err := model.GetDependencyTaskList(dependencyTaskId)
    if err != nil {
        logger.Errorf("获取依赖任务失败#主任务ID-%d#%s", taskModel.Id, err.Error())
        return
    }
    if len(tasks) == 0 {
        logger.Errorf("依赖任务列表为空#主任务ID-%d", taskModel.Id)
    }
    serviceTask := new(Task)
    for _, task := range tasks {
        task.Command = appendResultToCommand(task.Command, taskResult)
        task.Spec = fmt.Sprintf("依赖任务(主任务ID-%d)", taskModel.Id)
        serviceTask.Run(task)
    }
}

/**
 * 添加主任务执行结果到子任务命令中, 占位符{{.Code}} {{.Message}}
 */
func appendResultToCommand(command string, taskResult TaskResult) string {
    var code int8 = 0
    if taskResult.Err != nil {
        code = 1
    }
    data := map[string]interface{} {
        "Code": code,
        "Message": base64.StdEncoding.EncodeToString([]byte(taskResult.Result)),
    }
    var buf *bytes.Buffer = new(bytes.Buffer)
    tmpl, err := template.New("command").Parse(command)
    if err != nil {
        logger.Errorf("替换子任务命令占位符失败#%s", err.Error())
        return command
    }
    err = tmpl.Execute(buf, data)
    if err != nil {
        logger.Errorf("替换子任务命令占位符失败#%s", err.Error())
        return command
    }

    return buf.String()
}

// 发送任务结果通知
func SendNotification(taskModel models.Task, taskResult TaskResult)  {
    var statusName string
    // 未开启通知
    if taskModel.NotifyStatus == 0 {
        return
    }
    if taskModel.NotifyStatus == 1 && taskResult.Err == nil {
        // 执行失败才发送通知
        return
    }
    if taskModel.NotifyReceiverId == "" {
        return
    }
    if taskResult.Err != nil {
        statusName = "失败"
    } else {
        statusName = "成功"
    }
    // 发送通知
    msg := notify.Message{
        "task_type": taskModel.NotifyType,
        "task_receiver_id": taskModel.NotifyReceiverId,
        "name": taskModel.Name,
        "output": taskResult.Result,
        "status": statusName,
        "taskId": taskModel.Id,
    };
    notify.Push(msg)
}

// 执行具体任务
func execJob(handler Handler, taskModel models.Task) TaskResult  {
    defer func() {
       if err := recover(); err != nil {
           logger.Error("panic#service/task.go:execJob#", err)
       }
    } ()
    if taskModel.Multi == 0 {
        defer runInstance.done(taskModel.Id)
    }
    // 默认只运行任务一次
    var execTimes int8 = 1
    if (taskModel.RetryTimes > 0) {
        execTimes += taskModel.RetryTimes
    }
    var i int8 = 0
    var output string
    var err error
    for i < execTimes {
        output, err = handler.Run(taskModel)
        if err == nil {
            return TaskResult{Result: output, Err: err, RetryTimes: i}
        }
        i++
        if i < execTimes {
            logger.Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", taskModel.Id, i, output, err.Error())
            // 重试间隔时间，每次递增1分钟
            time.Sleep( time.Duration(i) * time.Minute)
        }
    }

    return TaskResult{Result: output, Err: err, RetryTimes: taskModel.RetryTimes}
}