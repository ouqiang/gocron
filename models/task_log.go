package models

import (
	"time"
)

// 任务执行日志
type TaskLog struct {
	Id        int       `xorm:"int pk autoincr"`
	TaskId    int       `xorm:"int not null"`                       // 任务ID
	StartTime time.Time `xorm:"datetime created"`                   // 开始执行时间
	EndTime   time.Time `xorm:"datetime updated"`                   // 执行完成（失败）时间
	Status    Status    `xorm:"tinyint notnull default 1"`          // 状态 1:执行中  2:执行完毕 0:执行失败
	Result    string    `xorm:"varchar(65535) notnull defalut '' "` // 执行结果
	Page      int       `xorm:"-"`
	PageSize  int       `xorm:"-"`
}

func (taskLog *TaskLog) Create() (insertId int, err error) {
	taskLog.Status = Running

	_, err = Db.Insert(taskLog)
	if err == nil {
		insertId = taskLog.Id
	}

	return
}

// 更新
func (taskLog *TaskLog) Update(id int, data CommonMap) (int64, error) {
	return Db.Table(taskLog).ID(id).Update(data)
}

func (taskLog *TaskLog) setStatus(id int, status Status) (int64, error) {
	return taskLog.Update(id, CommonMap{"status": status})
}

func (taskLog *TaskLog) List() ([]TaskLog, error) {
	taskLog.parsePageAndPageSize()
	list := make([]TaskLog, 0)
	err := Db.Desc("id").Limit(taskLog.PageSize, taskLog.pageLimitOffset()).Find(&list)

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
