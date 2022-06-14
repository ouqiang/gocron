package models

import "fmt"

type ProjectHost struct {
	Id        int `json:"id" xorm:"int pk autoincr"`
	ProjectId int `json:"project_id" xorm:"int not null index"`
	HostId    int `json:"host_id" xorm:"smallint not null index"`
}

func (ph *ProjectHost) Create() {
	exist, _ := Db.Exist(ph)
	if !exist {
		_, _ = Db.Insert(ph)
	}
}

func (ph *ProjectHost) RemoveForProject(project Project) (int64, error) {
	return Db.Where("project_id = ?", project.Id).Delete(ph)
}

func (ph *ProjectHost) GetHostsByProjectId(projectId int) ([]HostDetail, error) {
	list := make([]HostDetail, 0)
	err := Db.SQL(fmt.Sprintf("SELECT `ph`.`host_id`,`h`.`alias`,`h`.`name`,`h`.`port` FROM `%sproject_host` AS `ph` "+
		"LEFT JOIN `%shost` AS `h` ON `ph`.`host_id` = `h`.`id` "+
		"WHERE `ph`.`project_id` = %d", TablePrefix, TablePrefix, projectId)).Find(&list)
	//fmt.Println(rows, err)
	return list, err
}
