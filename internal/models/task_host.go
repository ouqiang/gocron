package models

type TaskHost struct {
	Id     int   `json:"id" xorm:"int pk autoincr"`
	TaskId int   `json:"task_id" xorm:"int not null index"`
	HostId int16 `json:"host_id" xorm:"smallint not null index"`
}

type TaskHostDetail struct {
	TaskHost `xorm:"extends"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Alias    string `json:"alias"`
}

func (TaskHostDetail) TableName() string {
	return TablePrefix + "task_host"
}

func hostTableName() []string {
	return []string{TablePrefix + "host", "h"}
}

func (th *TaskHost) Remove(taskId int) error {
	_, err := Db.Where("task_id = ?", taskId).Delete(new(TaskHost))

	return err
}

func (th *TaskHost) Add(taskId int, hostIds []int) error {

	err := th.Remove(taskId)
	if err != nil {
		return err
	}

	taskHosts := make([]TaskHost, len(hostIds))
	for i, value := range hostIds {
		taskHosts[i].TaskId = taskId
		taskHosts[i].HostId = int16(value)
	}

	_, err = Db.Insert(&taskHosts)

	return err
}

func (th *TaskHost) GetHostIdsByTaskId(taskId int) ([]TaskHostDetail, error) {
	list := make([]TaskHostDetail, 0)
	fields := "th.id,th.host_id,h.alias,h.name,h.port"
	err := Db.Alias("th").
		Join("LEFT", hostTableName(), "th.host_id=h.id").
		Where("th.task_id = ?", taskId).
		Cols(fields).
		Find(&list)

	return list, err
}

func (th *TaskHost) GetTaskIdsByHostId(hostId int16) ([]interface{}, error) {
	list := make([]TaskHost, 0)
	err := Db.Where("host_id = ?", hostId).Cols("task_id").Find(&list)
	if err != nil {
		return nil, err
	}

	taskIds := make([]interface{}, len(list))
	for i, value := range list {
		taskIds[i] = value.TaskId
	}

	return taskIds, err
}

// 判断主机id是否有引用
func (th *TaskHost) HostIdExist(hostId int16) (bool, error) {
	count, err := Db.Where("host_id = ?", hostId).Count(th)

	return count > 0, err
}
