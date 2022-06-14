package models

type ProcessHost struct {
	Id        int `json:"id" xorm:"int pk autoincr"`
	ProcessId int `json:"process_id" xorm:"int not null index"`
	HostId    int `json:"host_id" xorm:"smallint not null index"`
}

func (ph *ProcessHost) Create() {
	exist, _ := Db.Exist(ph)
	if !exist {
		_, _ = Db.Insert(ph)
	}
}

func (ph *ProcessHost) GetByProcess(process Process) []Host {
	hosts := make([]Host, 0)
	Db.Table(ProcessHost{}).Select("`host`.*").Join("LEFT", "host", "`host`.`id` = `process_host`.`host_id`").
		Where("process_id = ?", process.Id).Find(&hosts)
	return hosts
}

func (ph *ProcessHost) DeleteForProcess(process Process) {
	Db.Where("process_id = ?", process.Id).Delete(ph)
}
