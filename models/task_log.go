package models

import (
    "time"
)

type TaskType int8


// 任务执行日志
type TaskLog struct {
    Id        int64       `xorm:"bigint pk autoincr"`
    TaskId   int       `xorm:"int notnull index default 0"`       // 任务id
    Name     string    `xorm:"varchar(64) notnull"`               // 任务名称
    Spec     string    `xorm:"varchar(64) notnull"`               // crontab
    Protocol TaskProtocol  `xorm:"tinyint notnull"`               // 协议 1:http 2:ssh-command
    Command  string    `xorm:"varchar(512) notnull"`              // URL地址或shell命令
    Timeout  int       `xorm:"mediumint notnull default 0"`       // 任务执行超时时间(单位秒),0不限制
    Hostname string       `xorm:"varchar(512) notnull defalut '' "`   // SSH主机名，逗号分隔
    StartTime time.Time `xorm:"datetime created"`                   // 开始执行时间
    EndTime   time.Time `xorm:"datetime updated"`                   // 执行完成（失败）时间
    Status    Status    `xorm:"tinyint notnull default 1"`          // 状态 1:执行中  2:执行完毕 0:执行失败
    Result    string    `xorm:"varchar(65535) notnull defalut '' "` // 执行结果
    Page      int       `xorm:"-"`
    PageSize  int       `xorm:"-"`
}

func (taskLog *TaskLog) Create() (insertId int64, err error) {
    taskLog.Status = Running

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

func (taskLog *TaskLog) List() ([]TaskLog, error) {
    taskLog.parsePageAndPageSize()
    list := make([]TaskLog, 0)
    err := Db.Desc("id").Limit(taskLog.PageSize, taskLog.pageLimitOffset()).Find(&list)

    return list, err
}

// 清空表
func (TaskLog *TaskLog) Clear() (int64, error)  {
    return Db.Where("1=1").Delete(TaskLog);
}

func (taskLog *TaskLog) Total() (int64, error) {
    return Db.Count(taskLog)
}

func (taskLog *TaskLog) parsePageAndPageSize() {
    if taskLog.Page <= 0 {
        taskLog.Page = Page
    }
    if taskLog.PageSize >= 0 || taskLog.PageSize > MaxPageSize {
        taskLog.PageSize = PageSize
    }
}

func (taskLog *TaskLog) pageLimitOffset() int {
    return (taskLog.Page - 1) * taskLog.PageSize
}