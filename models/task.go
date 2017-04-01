package models

import (
	"time"
)

type Protocol int8

type TaskType int8

const (
	HTTP Protocol = 1
	SSHCommand Protocol  = 2
	SSHScript Protocol = 3
)

const (
	Timing TaskType = 1
	Delay  TaskType = 2
)

// 任务
type Task struct {
	Id int `xorm:"int pk autoincr"`
	Name string `xorm:"varchar(64) notnull"`          // 任务名称
	Spec string `xorm:"varchar(64) notnull"`          // crontab 时间格式
	Protocol Protocol `xorm:"tinyint notnull"`        // 协议 1:http 2:ssh-command 3:ssh-script
	Type TaskType `xorm:"tinyint notnull default 1"`      // 任务类型 1: 定时任务 2: 延时任务
	Command string `xorm:"varchar(512) notnull"`      // URL地址或shell命令
	Timeout int `xorm:"mediumint notnull default 0"`      // 任务执行超时时间(单位秒),0不限制
	Delay int `xorm:"int notnull default 0"`             // 延时任务，延时时间(单位秒)
	SshHosts string `xorm:"varchar(512) notnull defalut '' "` // SSH主机名, host id，逗号分隔
	Remark string `xorm:"varchar(512) notnull default ''"`    // 备注
	Created time.Time `xorm:"datetime notnull created"`       // 创建时间
	Deleted time.Time `xorm:"datetime deleted"`               // 删除时间
	Status Status `xorm:"tinyint notnull default 1"`    // 状态 1:正常 0:停止
	Page  int `xorm:"-"`
	PageSize int `xorm:"-"`
}

// 新增
func(task *Task) Create() (insertId int, err error) {
	task.Status = Enabled

	_, err =  Db.Insert(task)
	if err == nil {
		insertId = task.Id
	}

	return
}

// 更新
func(task *Task) Update(id int, data CommonMap) (int64, error) {
	return Db.Table(task).ID(id).Update(data)
}

// 删除
func(task *Task) Delete(id int) (int64, error) {
	return Db.Id(id).Delete(task)
}

// 禁用
func(task *Task) Disable(id int) (int64, error) {
	return task.Update(id, CommonMap{"status": Disabled})
}

// 激活
func(task *Task) Enable(id int) (int64, error)  {
	return task.Update(id, CommonMap{"status": Enabled})
}

func(task *Task) ActiveList() ([]Task, error)  {
	task.parsePageAndPageSize()
	list := make([]Task, 0)
	err := Db.Where("status = ?", Enabled).Desc("id").Find(&list)

	return list, err
}

func(task *Task) List() ([]Task, error) {
	task.parsePageAndPageSize()
	list := make([]Task, 0)
	err := Db.Desc("id").Limit(task.PageSize, task.pageLimitOffset()).Find(&list)

	return list, err
}

func(taskLog *TaskLog) Total() (int64, error) {
	return Db.Count(taskLog)
}

func(taskLog *TaskLog) parsePageAndPageSize()  {
	if (taskLog.Page <= 0) {
		taskLog.Page = Page
	}
	if (taskLog.PageSize >= 0 || taskLog.PageSize > MaxPageSize) {
		taskLog.PageSize = PageSize
	}
}

func(taskLog *TaskLog) pageLimitOffset() int  {
	return (taskLog.Page - 1) * taskLog.PageSize
}