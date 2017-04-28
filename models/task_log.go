package models

import (
    "time"
    "github.com/go-xorm/xorm"
)

type TaskType int8


// 任务执行日志
type TaskLog struct {
    Id        int64       `xorm:"bigint pk autoincr"`
    TaskId   int       `xorm:"int notnull index default 0"`       // 任务id
    Name     string    `xorm:"varchar(32) notnull"`               // 任务名称
    Spec     string    `xorm:"varchar(64) notnull"`               // crontab
    Protocol TaskProtocol  `xorm:"tinyint notnull"`               // 协议 1:http 2:ssh-command 3:系统命令
    Command  string    `xorm:"varchar(256) notnull"`              // URL地址或shell命令
    Timeout  int       `xorm:"mediumint notnull default 0"`       // 任务执行超时时间(单位秒),0不限制
    RetryTimes int8    `xorm:"tinyint notnull default 0"`           // 任务重试次数
    Hostname string       `xorm:"varchar(128) notnull defalut '' "`   // SSH主机名，逗号分隔
    StartTime time.Time `xorm:"datetime created"`                   // 开始执行时间
    EndTime   time.Time `xorm:"datetime updated"`                   // 执行完成（失败）时间
    Status    Status    `xorm:"tinyint notnull default 1"`          // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成)
    Result    string    `xorm:"mediumtext notnull defalut '' "` // 执行结果
    TotalTime int       `xorm:"-"` // 执行总时长
    BaseModel   `xorm:"-"`
}

func (taskLog *TaskLog) Create() (insertId int64, err error) {
    _, err = Db.Insert(taskLog)
    if err == nil {
        insertId = taskLog.Id
    }

    return
}

// 更新
func (taskLog *TaskLog) Update(id int64, data CommonMap) (int64, error) {
    return Db.Table(taskLog).ID(id).Update(data)
}

func (taskLog *TaskLog) setStatus(id int64, status Status) (int64, error) {
    return taskLog.Update(id, CommonMap{"status": status})
}

func (taskLog *TaskLog) List(params CommonMap) ([]TaskLog, error) {
    taskLog.parsePageAndPageSize(params)
    list := make([]TaskLog, 0)
    session := Db.Desc("id")
    taskLog.parseWhere(session, params)
    err := session.Limit(taskLog.PageSize, taskLog.pageLimitOffset()).Find(&list)
    if len(list) > 0 {
        for i, item := range list {
            endTime := item.EndTime
            if item.Status == Running {
                endTime = time.Now()
            }
            execSeconds := endTime.Sub(item.StartTime).Seconds()
            list[i].TotalTime = int(execSeconds)
        }
    }

    return list, err
}

// 清空表
func (TaskLog *TaskLog) Clear() (int64, error)  {
    return Db.Where("1=1").Delete(TaskLog);
}

func (taskLog *TaskLog) Total(params CommonMap) (int64, error) {
    session := Db.NewSession()
    taskLog.parseWhere(session, params)
    return session.Count(taskLog)
}

// 解析where
func (taskLog *TaskLog) parseWhere(session *xorm.Session, params CommonMap)  {
    if len(params) == 0 {
        return
    }
    taskId, ok := params["TaskId"]
    if ok && taskId.(int) > 0 {
        session.And("task_id = ?", taskId)
    }
    protocol, ok := params["Protocol"]
    if ok && protocol.(int) > 0 {
        session.And("protocol = ?", protocol)
    }
    status, ok := params["Status"]
    if ok && status.(int) > -1 {
        session.And("status = ?", status)
    }
}