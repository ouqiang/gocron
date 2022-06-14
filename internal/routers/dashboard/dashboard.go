package dashboard

import (
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"gopkg.in/macaron.v1"
)

type Dashboard struct {
	TotalGroup   map[string]int64          `json:"totalGroup"`
	ActiveUsers  []models.ActiveUser       `json:"activeUsers"`
	ProjectTasks []models.ProjectTaskChart `json:"projectTasks"`
	//NewData         map[string]map[string]int
}

func Index(ctx *macaron.Context) string {
	data := Dashboard{TotalGroup: map[string]int64{}}

	t := models.Task{}
	taskCount, _ := t.Total(models.CommonMap{})
	data.TotalGroup["taskCount"] = taskCount

	p := models.Process{}
	processCount, _ := p.Total(models.CommonMap{})
	data.TotalGroup["processCount"] = processCount
	u := models.User{}
	userCount, _ := u.Total()
	data.TotalGroup["userCount"] = userCount

	pro := models.Project{}
	projectCount, _ := pro.Total(models.CommonMap{})
	data.TotalGroup["projectCount"] = projectCount

	log := models.OperateLog{}
	users, _ := log.GetActiveUsers()
	data.ActiveUsers = users

	project := models.Project{}
	data.ProjectTasks = project.GetProjectTasksChart()

	return utils.JsonResp.Success(utils.SuccessContent, data)
}
