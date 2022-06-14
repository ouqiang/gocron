package models

import "github.com/go-xorm/xorm"

type ProcessStatus int8

const (
	ProcessStart ProcessStatus = 1
	ProcessStop  ProcessStatus = 2
)

type Process struct {
	Id          int             `json:"id" xorm:"int pk autoincr"`
	Name        string          `json:"name" xorm:"varchar(256) notnull"`
	Command     string          `json:"command" xorm:"varchar(256) notnull"`
	Status      ProcessStatus   `json:"status" xorm:"tinyint notnull default 2"`
	Tag         string          `json:"tag" xorm:"varchar(32) notnull default ''"`
	NumProc     int             `json:"num_proc" xorm:"tinyint notnull default 1"`
	AutoStart   int8            `json:"auto_start" xorm:"tinyint notnull default 0"`
	AuthRestart int8            `json:"auth_restart" xorm:"tinyint notnull default 0"`
	Enable      bool            `json:"enable" xorm:"tinyint notnull default 1"`
	LogFile     string          `json:"log_file" xorm:"varchar(256) notnull default"`
	Workers     []ProcessWorker `json:"workers" xorm:"-"`
	Hosts       []Host          `json:"hosts" xorm:"-"`
	BaseModel   `json:"-" xorm:"-"`
}

func (p *Process) Create() (processId int, err error) {
	_, err = Db.Insert(p)
	if err == nil {
		processId = p.Id
	}
	return
}

func (p *Process) UpdateBean(id int) (int64, error) {
	return Db.ID(id).
		Cols(`name,command,tag,num_proc,auto_start,auto_restart,log_file`).
		Update(p)
}

func (p *Process) Update(id int, data CommonMap) (int64, error) {
	return Db.Table(p).ID(id).Update(data)
}

func (p *Process) Get(id int) error {
	_, err := Db.Where("id=?", id).Get(p)
	return err
}

func (p *Process) Total(params CommonMap) (int64, error) {
	session := Db.Alias("p")
	p.parseWhere(session, params)
	return session.Count(p)
}

func (p *Process) List(params CommonMap) ([]Process, error) {
	session := Db.Alias("p")
	processList := make([]Process, 0)
	p.parsePageAndPageSize(params)
	p.parseWhere(session, params)

	err := session.Limit(p.PageSize, p.pageLimitOffset()).Find(&processList)
	if err != nil {
		return nil, err
	}

	return processList, nil
}

func (p *Process) GetStartWorkerTotal() (int64, error) {
	return Db.Where("process_id = ? AND state = ?", p.Id, Start).Count(ProcessWorker{})
}

// 解析where
func (p *Process) parseWhere(session *xorm.Session, params CommonMap) {
	if len(params) == 0 {
		return
	}
	id, ok := params["Id"]
	if ok && id.(int) > 0 {
		session.And("p.id = ?", id)
	}
	name, ok := params["Name"]
	if ok && name.(string) != "" {
		session.And("p.name LIKE ?", "%"+name.(string)+"%")
	}
	status, ok := params["Status"]
	if ok && status.(int) > -1 {
		session.And("p.status = ?", status)
	}

	command, ok := params["Command"]
	if ok && command.(string) != "" {
		session.And("p.command LIKE ?", "%"+command.(string)+"%")
	}
}
