package models

import (
    "time"
    "github.com/ouqiang/gocron/modules/ssh"
)

type TaskProtocol int8

const (
    TaskHTTP TaskProtocol = iota + 1 // HTTP协议
    TaskSSH  // SSH命令
    TaskLocalCommand // 本地命令
)

// 任务
type Task struct {
    Id       int       `xorm:"int pk autoincr"`
    Name     string    `xorm:"varchar(64) notnull"`              // 任务名称
    Spec     string    `xorm:"varchar(64) notnull"`              // crontab
    Protocol TaskProtocol  `xorm:"tinyint notnull"`              // 协议 1:http 2:ssh-command 3: 本地命令
    Command  string    `xorm:"varchar(512) notnull"`             // URL地址或shell命令
    Timeout  int       `xorm:"mediumint notnull default 0"`      // 任务执行超时时间(单位秒),0不限制
    HostId   int16    `xorm:"smallint notnull default 0"`        // SSH host id，
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
    Alias string
    AuthType ssh.HostAuthType
    PrivateKey string
}

func (TaskHost) TableName() string  {
    return "task"
}

// 新增
func (task *Task) Create() (insertId int, err error) {
    _, err = Db.Insert(task)
    if err == nil {
        insertId = task.Id
    }

    return
}

func (task *Task) UpdateBean(id int) (int64, error)  {
    return Db.ID(id).UseBool("status").Update(task)
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

// 获取所有激活任务
func (task *Task) ActiveList() ([]TaskHost, error) {
    task.parsePageAndPageSize()
    list := make([]TaskHost, 0)
    fields := "t.*, host.alias,host.name,host.username,host.password,host.port,host.auth_type,host.private_key"
    err := Db.Alias("t").Join("LEFT", "host", "t.host_id=host.id").Where("t.status = ?", Enabled).Cols(fields).Find(&list)

    return list, err
}

// 获取某个主机下的所有激活任务
func (task *Task) ActiveListByHostId(hostId int16) ([]TaskHost, error)  {
    task.parsePageAndPageSize()
    list := make([]TaskHost, 0)
    fields := "t.*, host.alias,host.name,host.username,host.password,host.port,host.auth_type,host.private_key"
    err := Db.Alias("t").Join("LEFT", "host", "t.host_id=host.id").Where("t.status = ? AND t.host_id = ?", Enabled, hostId).Cols(fields).Find(&list)

    return list, err
}

// 判断主机id是否有引用
func (task *Task) HostIdExist(hostId int16) (bool, error) {
    count, err := Db.Where("host_id = ?", hostId).Count(task);

    return count > 0, err
}

// 判断任务名称是否存在
func (task *Task) NameExist(name string, id int) (bool, error)  {
    if id > 0 {
        count, err := Db.Where("name = ? AND status = ? AND id != ?", name, Enabled, id).Count(task);
        return count > 0, err
    }
    count, err := Db.Where("name = ? AND status = ?", name, Enabled).Count(task);

    return count > 0, err
}

func(task *Task) Detail(id int) (TaskHost, error)  {
    taskHost := TaskHost{}
    fields := "t.*, host.alias,host.name,host.username,host.password,host.port,host.auth_type,host.private_key"
    _, err := Db.Alias("t").Join("LEFT", "host", "t.host_id=host.id").Where("t.id=?", id).Cols(fields).Get(&taskHost)

    return taskHost, err
}

func (task *Task) List() ([]TaskHost, error) {
    task.parsePageAndPageSize()
    list := make([]TaskHost, 0)
    fields := "t.*, host.alias"
    err := Db.Alias("t").Join("LEFT", "host", "t.host_id=host.id").Cols(fields).Desc("t.id").Limit(task.PageSize, task.pageLimitOffset()).Find(&list)

    return list, err
}

func (task *Task) Total() (int64, error) {
    return Db.Count(task)
}

func (task *Task) parsePageAndPageSize() {
    if task.Page <= 0 {
        task.Page = Page
    }
    if task.PageSize >= 0 || task.PageSize > MaxPageSize {
        task.PageSize = PageSize
    }
}

func (task *Task) pageLimitOffset() int {
    return (task.Page - 1) * task.PageSize
}