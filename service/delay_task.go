package service

import (
    "github.com/ouqiang/gocron/models"
    "time"
    "github.com/ouqiang/gocron/modules/logger"
    "math"
    "github.com/ouqiang/gocron/modules/httpclient"
    "strings"
    "github.com/ouqiang/timewheel"
)

var tw  *timewheel.TimeWheel

type DelayTask struct {}

// 从数据库中取出所有延迟任务
func (task *DelayTask) Initialize(tick time.Duration, slots int)  {
    tw = timewheel.New(tick, slots, task.Run)
    tw.Start()
    taskModel := new(models.DelayTask)
    currentTime := time.Now()
    taskNum, err := taskModel.ActiveNum(currentTime)
    if err != nil {
        logger.Error("延迟任务初始化#获取待执行的任务失败", err)
        return
    }
    if taskNum == 0 {
        logger.Debugf("延迟任务初始化#待执行的任务数量为0")
        return
    }
    pageSize := 100
    totalPage := int( math.Ceil(float64(taskNum) / float64(pageSize)) )
    logger.Infof("延迟任务初始化#待执行的任务数量-%d#共%d页#每页取%d条", taskNum, totalPage, pageSize)
    taskModel.PageSize = pageSize
    for page := 1; page <= totalPage; page++ {
        taskModel.Page = page
        logger.Debugf("延迟任务初始化#取出任务列表#第%d页", page)
        taskList, err := taskModel.ActiveList(currentTime)
        if err != nil {
            logger.Error("延迟任务初始化#获取任务列表失败", err)
        }
        task.BatchAdd(taskList)
    }
    logger.Info("延迟任务初始化完成")
}

// 批量添加任务
func (task *DelayTask) BatchAdd(taskList []models.DelayTask)  {
    for _, item := range(taskList) {
        task.Add(item)
    }
}

// 添加任务
func (task *DelayTask) Add(taskModel models.DelayTask)  {
    currentTimestamp := time.Now().Unix()
    execTimestamp := taskModel.Created.Unix() + int64(taskModel.Delay)
    // 时间过期, 立即执行任务
    data := []interface{}{taskModel.Id, taskModel.Url, taskModel.Params}
    if execTimestamp <= currentTimestamp {
        go task.Run(data)
        return
    }
    delay := execTimestamp - currentTimestamp
    tw.Add(time.Duration(delay) * time.Second, data)
}

// 运行任务
func (task *DelayTask) Run(data []interface{})  {
    if len(data) < 3 {
        logger.Errorf("延时任务开始执行#参数不足#%+v", data)
        return
    }
    id := data[0].(int64)
    url := data[1].(string)
    params := data[2].(string)
    if id <= 0 || url == "" {
        logger.Errorf("延时任务开始执行#参数为空#%+v", data)
        return
    }
    taskModel := new(models.DelayTask)
    _, err := taskModel.UpdateStatus(id, models.Running)
    if err != nil {
        logger.Error("延迟任务开始执行#更新任务状态失败", err)
        return
    }
    timeout := 300
    tryTimes := 3
    success := false
    logger.Infof("延迟任务开始执行#id-%d#url-%s#params-%s", id, url, params)
    for i := 0; i < tryTimes; {
        response := httpclient.PostParams(url, params, timeout)
        if response.StatusCode == 200 && strings.TrimSpace(response.Body) == "success"{
            success = true
            break;
        }
        i++
        if i < tryTimes {
            logger.Errorf("延迟任务执行失败#重试第%d次#任务Id-%d#HTTP状态码-%d#HTTP-BODY-%s",
            i,id,response.StatusCode,response.Body)
            time.Sleep(30 * time.Second)
        }
    }
    logger.Infof("延迟任务执行完成#id-%d", id)
    var status models.Status
    if success {
        status = models.Finish
    } else {
        status = models.Failure
    }
    _ ,err = taskModel.UpdateStatus(id, status)
    if err != nil {
        logger.Error("延迟任务执行完成#更新任务状态失败", err)
    }
}

func (task *DelayTask) Stop()  {
    tw.Stop()
}