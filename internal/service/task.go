package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ouqiang/goutil"

	"github.com/jakecoffman/cron"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/httpclient"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/notify"
	rpcClient "github.com/ouqiang/gocron/internal/modules/rpc/client"
	pb "github.com/ouqiang/gocron/internal/modules/rpc/proto"
)

var (
	ServiceTask Task
)

var (
	// 定时任务调度管理器
	serviceCron *cron.Cron

	// 同一任务是否有实例处于运行中
	runInstance Instance

	// 任务计数-正在运行的任务
	taskCount TaskCount

	// 并发队列, 限制同时运行的任务数量
	concurrencyQueue ConcurrencyQueue
)

// 并发队列
type ConcurrencyQueue struct {
	queue chan struct{}
}

func (cq *ConcurrencyQueue) Add() {
	cq.queue <- struct{}{}
}

func (cq *ConcurrencyQueue) Done() {
	<-cq.queue
}

// 任务计数
type TaskCount struct {
	wg   sync.WaitGroup
	exit chan struct{}
}

func (tc *TaskCount) Add() {
	tc.wg.Add(1)
}

func (tc *TaskCount) Done() {
	tc.wg.Done()
}

func (tc *TaskCount) Exit() {
	tc.wg.Done()
	<-tc.exit
}

func (tc *TaskCount) Wait() {
	tc.Add()
	tc.wg.Wait()
	close(tc.exit)
}

// 任务ID作为Key
type Instance struct {
	m sync.Map
}

// 是否有任务处于运行中
func (i *Instance) has(key int) bool {
	_, ok := i.m.Load(key)

	return ok
}

func (i *Instance) add(key int) {
	i.m.Store(key, struct{}{})
}

func (i *Instance) done(key int) {
	i.m.Delete(key)
}

type Task struct{}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

// 初始化任务, 从数据库取出所有任务, 添加到定时任务并运行
func (task Task) Initialize() {
	serviceCron = cron.New()
	serviceCron.Start()
	concurrencyQueue = ConcurrencyQueue{queue: make(chan struct{}, app.Setting.ConcurrencyQueue)}
	taskCount = TaskCount{sync.WaitGroup{}, make(chan struct{})}
	go taskCount.Wait()

	logger.Info("开始初始化定时任务")
	taskModel := new(models.Task)
	taskNum := 0
	page := 1
	pageSize := 1000
	maxPage := 1000
	for page < maxPage {
		taskList, err := taskModel.ActiveList(page, pageSize)
		if err != nil {
			logger.Fatalf("定时任务初始化#获取任务列表错误: %s", err)
		}
		if len(taskList) == 0 {
			break
		}
		for _, item := range taskList {
			task.Add(item)
			taskNum++
		}
		page++
	}
	logger.Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

// 批量添加任务
func (task Task) BatchAdd(tasks []models.Task) {
	for _, item := range tasks {
		task.RemoveAndAdd(item)
	}
}

// 删除任务后添加
func (task Task) RemoveAndAdd(taskModel models.Task) {
	task.Remove(taskModel.Id)
	task.Add(taskModel)
}

// 添加任务
func (task Task) Add(taskModel models.Task) {
	if taskModel.Level == models.TaskLevelChild {
		logger.Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", taskModel.Id)
		return
	}
	taskFunc := createJob(taskModel)
	if taskFunc == nil {
		logger.Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
		return
	}

	cronName := strconv.Itoa(taskModel.Id)
	err := goutil.PanicToError(func() {
		serviceCron.AddFunc(taskModel.Spec, taskFunc, cronName)
	})
	if err != nil {
		logger.Error("添加任务到调度器失败#", err)
	}
}

func (task Task) NextRunTime(taskModel models.Task) time.Time {
	if taskModel.Level != models.TaskLevelParent ||
		taskModel.Status != models.Enabled {
		return time.Time{}
	}
	entries := serviceCron.Entries()
	taskName := strconv.Itoa(taskModel.Id)
	for _, item := range entries {
		if item.Name == taskName {
			return item.Next
		}
	}

	return time.Time{}
}

// 停止运行中的任务
func (task Task) Stop(ip string, port int, id int64) {
	rpcClient.Stop(ip, port, id)
}

func (task Task) Remove(id int) {
	serviceCron.RemoveJob(strconv.Itoa(id))
}

// 等待所有任务结束后退出
func (task Task) WaitAndExit() {
	serviceCron.Stop()
	taskCount.Exit()
}

// 直接运行任务
func (task Task) Run(taskModel models.Task) {
	go createJob(taskModel)()
}

type Handler interface {
	Run(taskModel models.Task, taskUniqueId int64) (string, error)
}

// HTTP任务
type HTTPHandler struct{}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(taskModel models.Task, taskUniqueId int64) (result string, err error) {
	if taskModel.Timeout <= 0 || taskModel.Timeout > HttpExecTimeout {
		taskModel.Timeout = HttpExecTimeout
	}
	var resp httpclient.ResponseWrapper
	if taskModel.HttpMethod == models.TaskHTTPMethodGet {
		resp = httpclient.Get(taskModel.Command, taskModel.Timeout)
	} else {
		urlFields := strings.Split(taskModel.Command, "?")
		taskModel.Command = urlFields[0]
		var params string
		if len(urlFields) >= 2 {
			params = urlFields[1]
		}
		resp = httpclient.PostParams(taskModel.Command, params, taskModel.Timeout)
	}
	// 返回状态码非200，均为失败
	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
	}

	return resp.Body, err
}

// RPC调用执行任务
type RPCHandler struct{}

func (h *RPCHandler) Run(taskModel models.Task, taskUniqueId int64) (result string, err error) {
	taskRequest := new(pb.TaskRequest)
	taskRequest.Timeout = int32(taskModel.Timeout)
	taskRequest.Command = taskModel.Command
	taskRequest.Id = taskUniqueId
	resultChan := make(chan TaskResult, len(taskModel.Hosts))
	for _, taskHost := range taskModel.Hosts {
		go func(th models.TaskHostDetail) {
			output, err := rpcClient.Exec(th.Name, th.Port, taskRequest)
			errorMessage := ""
			if err != nil {
				errorMessage = err.Error()
			}
			outputMessage := fmt.Sprintf("主机: [%s-%s:%d]\n%s\n%s\n\n",
				th.Alias, th.Name, th.Port, errorMessage, output,
			)
			resultChan <- TaskResult{Err: err, Result: outputMessage}
		}(taskHost)
	}

	var aggregationErr error = nil
	aggregationResult := ""
	for i := 0; i < len(taskModel.Hosts); i++ {
		taskResult := <-resultChan
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
		aggregationHost := ""
		for _, host := range taskModel.Hosts {
			aggregationHost += fmt.Sprintf("%s - %s<br>", host.Alias, host.Name)
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
	result := taskResult.Result
	if taskResult.Err != nil {
		status = models.Failure
	} else {
		status = models.Finish
	}
	return taskLogModel.Update(taskLogId, models.CommonMap{
		"retry_times": taskResult.RetryTimes,
		"status":      status,
		"result":      result,
	})

}

func createJob(taskModel models.Task) cron.FuncJob {
	handler := createHandler(taskModel)
	if handler == nil {
		return nil
	}
	taskFunc := func() {
		taskCount.Add()
		defer taskCount.Done()

		taskLogId := beforeExecJob(taskModel)
		if taskLogId <= 0 {
			return
		}

		if taskModel.Multi == 0 {
			runInstance.add(taskModel.Id)
			defer runInstance.done(taskModel.Id)
		}

		concurrencyQueue.Add()
		defer concurrencyQueue.Done()

		logger.Infof("开始执行任务#%s#命令-%s", taskModel.Name, taskModel.Command)
		taskResult := execJob(handler, taskModel, taskLogId)
		logger.Infof("任务完成#%s#命令-%s", taskModel.Name, taskModel.Command)
		afterExecJob(taskModel, taskResult, taskLogId)
	}

	return taskFunc
}

func createHandler(taskModel models.Task) Handler {
	var handler Handler = nil
	switch taskModel.Protocol {
	case models.TaskHTTP:
		handler = new(HTTPHandler)
	case models.TaskRPC:
		handler = new(RPCHandler)
	}

	return handler
}

// 任务前置操作
func beforeExecJob(taskModel models.Task) (taskLogId int64) {
	if taskModel.Multi == 0 && runInstance.has(taskModel.Id) {
		createTaskLog(taskModel, models.Cancel)
		return
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
func afterExecJob(taskModel models.Task, taskResult TaskResult, taskLogId int64) {
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
func execDependencyTask(taskModel models.Task, taskResult TaskResult) {
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
	tasks, err := model.GetDependencyTaskList(dependencyTaskId)
	if err != nil {
		logger.Errorf("获取依赖任务失败#主任务ID-%d#%s", taskModel.Id, err.Error())
		return
	}
	if len(tasks) == 0 {
		logger.Errorf("依赖任务列表为空#主任务ID-%d", taskModel.Id)
	}
	for _, task := range tasks {
		task.Spec = fmt.Sprintf("依赖任务(主任务ID-%d)", taskModel.Id)
		ServiceTask.Run(task)
	}
}

// 发送任务结果通知
func SendNotification(taskModel models.Task, taskResult TaskResult) {
	var statusName string
	// 未开启通知
	if taskModel.NotifyStatus == 0 {
		return
	}
	if taskModel.NotifyStatus == 3 {
		// 关键字匹配通知
		if !strings.Contains(taskResult.Result, taskModel.NotifyKeyword) {
			return
		}
	}
	if taskModel.NotifyStatus == 1 && taskResult.Err == nil {
		// 执行失败才发送通知
		return
	}
	if taskModel.NotifyType != 3 && taskModel.NotifyReceiverId == "" {
		return
	}
	if taskResult.Err != nil {
		statusName = "失败"
	} else {
		statusName = "成功"
	}
	// 发送通知
	msg := notify.Message{
		"task_type":        taskModel.NotifyType,
		"task_receiver_id": taskModel.NotifyReceiverId,
		"name":             taskModel.Name,
		"output":           taskResult.Result,
		"status":           statusName,
		"task_id":          taskModel.Id,
	}
	notify.Push(msg)
}

// 执行具体任务
func execJob(handler Handler, taskModel models.Task, taskUniqueId int64) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("panic#service/task.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var execTimes int8 = 1
	if taskModel.RetryTimes > 0 {
		execTimes += taskModel.RetryTimes
	}
	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		output, err = handler.Run(taskModel, taskUniqueId)
		if err == nil {
			return TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		if i < execTimes {
			logger.Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", taskModel.Id, i, output, err.Error())
			if taskModel.RetryInterval > 0 {
				time.Sleep(time.Duration(taskModel.RetryInterval) * time.Second)
			} else {
				// 默认重试间隔时间，每次递增1分钟
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}

	return TaskResult{Result: output, Err: err, RetryTimes: taskModel.RetryTimes}
}
