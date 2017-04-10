package models

import (
    "time"
)

type Protocol int8

type TaskType int8

const (
    HTTP       Protocol = iota + 1 // HTTP协议
    SSHCommand  // SSHM命令
)

const (
    Timing TaskType = iota + 1 // 定时任务
    Delay // 延时任务
)

// 任务
type Task struct {
    Id       int       `xorm:"int pk autoincr"`
    Name     string    `xorm:"varchar(64) notnull"`              // 任务名称
    Spec     string    `xorm:"varchar(64) notnull"`              // crontab
    Protocol Protocol  `xorm:"tinyint notnull"`                  // 协议 1:http 2:ssh-command
    Type     TaskType  `xorm:"tinyint notnull default 1"`        // 任务类型 1: 定时任务 2: 延时任务
    Command  string    `xorm:"varchar(512) notnull"`             // URL地址或shell命令
    Timeout  int       `xorm:"mediumint notnull default 0"`      // 任务执行超时时间(单位秒),0不限制
    Delay    int       `xorm:"int notnull default 0"`            // 延时任务，延时时间(单位秒)
    HostId   int16    `xorm:"smallint notnull default 0"`       // SSH host id，
    Remark   string    `xorm:"varchar(512) notnull default ''"`  // 备注
    Created  time.Time `xorm:"datetime notnull created"`         // 创建时间
    Deleted  time.Time `xorm:"datetime deleted"`                 // 删除时间
    Status   Status    `xorm:"tinyint notnull default 1"`        // 状态 1:正常 0:停止
    Page     int       `xorm:"-"`
    PageSize int       `xorm:"-"`
}

type TaskHost struct {
    Task `xorm:"extends"`
    Name string
    Port int
    Username string
    Password string
}

func (TaskHost) TableName() string  {
    return "task"
}

// 新增
func (task *Task) Create() (insertId int, err error) {
    task.Status = Enabled

    _, err = Db.Insert(task)
    if err == nil {
        insertId = task.Id
    }

    return
}

// 更新
func (task *Task) Update(id int, data CommonMap) (int64, error) {
    return Db.Table(task).ID(id).Update(data)
}

// 删除
func (task *Task) Delete(id int) (int64, error) {
    return Db.Id(id).Delete(task)
}

// 禁用
func (task *Task) Disable(id int) (int64, error) {
    return task.Update(id, CommonMap{"status": Disabled})
}

// 激活
func (task *Task) Enable(id int) (int64, error) {
    return task.Update(id, CommonMap{"status": Enabled})
}

func (task *Task) ActiveList() ([]TaskHost, error) {
    task.parsePageAndPageSize()
    list := make([]TaskHost, 0)
    fields := "t.*, host.name,host.username,host.password,host.port"
    err := Db.Alias("t").Join("LEFT", "host", "t.host_id=host.id").Where("status = ?", Enabled).Cols(fields).Find(&list)

    return list, err
}

func(task *Task) Detail(id int) error  {
    list := make([]TaskHost, 0)
    fields := "t.*, host.name,host.username,host.password,host.port"
    err := Db.Alias("t").Join("LEFT", "host", "t.host_id=host.id").Cols(fields).Find(list)

    return err
}

func (task *Task) List() ([]TaskHost, error) {
    task.parsePageAndPageSize()
    list := make([]TaskHost, 0)
    fields := "t.*, host.name"
    err := Db.Alias("t").Join("LEFT", "host", "t.host_id=host.id").Cols(fields).Desc("t.id").Limit(task.PageSize, task.pageLimitOffset()).Find(&list)

    return list, err
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
