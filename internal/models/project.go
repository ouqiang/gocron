package models

import (
	"github.com/go-xorm/xorm"
	"time"
)

type Project struct {
	Id           int       `json:"id" xorm:"int pk autoincr"`
	Name         string    `json:"name" xorm:"varchar(256) notnull"`
	Code         string    `json:"code" xorm:"varchar(16) notnull"`
	Remark       string    `json:"remark" xorm:"varchar(256) notnull default '' "` // 备注
	CreatedAt    time.Time `json:"created_at" xorm:"datetime notnull created"`
	CreateUserId int       `json:"create_user_id" xorm:"int notnull default 0"`
	UpdatedAt    time.Time `json:"updated_at" xorm:"datetime notnull updated"`
	UpdateUserId int       `json:"update_user_id" xorm:"int notnull default 0"`
	Hosts        []Host    `json:"hosts" xorm:"-"`
	BaseModel    `json:"-" xorm:"-"`
}

type ProjectTaskChart struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func (p *Project) Create() (int64, error) {
	return Db.Insert(p)
}

func (p *Project) Update() (int64, error) {
	return Db.Cols(`name,code,remark,update_user_id`).Where("id = ?", p.Id).Update(p)
}

func (p *Project) Total(params CommonMap) (int64, error) {
	session := Db.Alias("p")
	p.parseWhere(session, params)
	return session.Count(p)
}

func (p *Project) List(params CommonMap) ([]Project, error) {
	session := Db.Alias("p")
	projects := make([]Project, 0)
	p.parsePageAndPageSize(params)
	p.parseWhere(session, params)

	err := session.Limit(p.PageSize, p.pageLimitOffset()).Find(&projects)
	if err != nil {
		return nil, err
	}

	return p.setHostsForProjects(projects), nil
}

func (p *Project) setHostsForProjects(projects []Project) []Project {
	var projectIds []int
	for _, project := range projects {
		projectIds = append(projectIds, project.Id)
	}
	if len(projectIds) == 0 {
		return projects
	}
	phList := make([]ProjectHost, 0)
	_ = Db.In("project_id", projectIds).Find(&phList)

	var hostIds []int
	var projectHostGroup = make(map[int][]int)
	for _, ph := range phList {
		hostIds = append(hostIds, ph.HostId)
		values, ok := projectHostGroup[ph.ProjectId]
		if ok {
			projectHostGroup[ph.ProjectId] = append(values, ph.HostId)
		} else {
			projectHostGroup[ph.ProjectId] = []int{ph.HostId}
		}
	}
	hosts := make([]Host, 0)
	_ = Db.In("id", hostIds).Find(&hosts)
	var hostGroup = make(map[int]Host)
	for _, host := range hosts {
		hostGroup[host.Id] = host
	}

	for i, project := range projects {
		hostIds := projectHostGroup[project.Id]
		for _, hostId := range hostIds {
			projects[i].Hosts = append(projects[i].Hosts, hostGroup[hostId])
		}
	}
	return projects
}

// 解析where
func (p *Project) parseWhere(session *xorm.Session, params CommonMap) {
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
}

func (p *Project) All() ([]Project, interface{}) {
	projects := make([]Project, 0)
	err := Db.Find(&projects)
	return projects, err
}

func (p *Project) GetProjectTasksChart() []ProjectTaskChart {
	result := make([]ProjectTaskChart, 0)
	_ = Db.SQL("SELECT `p`.`name`,count(`t`.`id`) AS value FROM `project` AS `p` LEFT JOIN `task` AS `t` ON `t`.`project_id` = `p`.`id` GROUP BY `p`.`id`").Find(&result)
	return result
}
