package models

import "time"

type WorkerState int8

const (
	Pending  WorkerState = 0 //挂起
	Starting WorkerState = 1 // 启动中
	Start    WorkerState = 2 //已启动
	Stopping WorkerState = 3 //停止中
	Stop     WorkerState = 4 //已停止
	Exited   WorkerState = 5 //异常退出
)

type ProcessWorker struct {
	Id        int64       `json:"id" xorm:"int pk autoincr"`
	HostId    int         `json:"host_id" xorm:"int host_id"`
	ProcessId int         `json:"process_id" xorm:"int process_id"`
	Pid       int64       `json:"pid" xorm:"int pid"`
	IsValid   int8        `json:"is_valid" xorm:"tinyint notnull default 0"`
	State     WorkerState `json:"state" xorm:"tinyint notnull default 0"`
	StartAt   time.Time   `json:"start_at" xorm:"datetime notnull created"`
	BaseModel `json:"-" xorm:"-"`
}

func (pw *ProcessWorker) Create() (err error) {
	_, err = Db.Insert(pw)
	return err
}

func (pw *ProcessWorker) GetByProcess(process Process) ([]ProcessWorker, error) {
	var workers []ProcessWorker
	err := Db.Where("process_id = ? AND is_valid = ?", process.Id, 1).Find(&workers)
	return workers, err
}

func (pw *ProcessWorker) GetLimitByProcess(process Process, limit int) ([]ProcessWorker, error) {
	var workers []ProcessWorker
	err := Db.Where("process_id = ? AND is_valid = ?", process.Id, 1).Limit(limit).Find(&workers)
	return workers, err
}

func (pw *ProcessWorker) Update() error {
	_, err := Db.Where("id = ?", pw.Id).Cols(`host_id,process_id,pid,state,start_at,is_valid`).Update(pw)
	return err
}

func (pw *ProcessWorker) SetState(state WorkerState) error {
	_, err := Db.Where("id = ?", pw.Id).Update(ProcessWorker{State: state})
	return err
}
