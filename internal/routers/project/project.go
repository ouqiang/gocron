package project

import (
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/gocron/internal/routers/base"
	"github.com/ouqiang/gocron/internal/service"
	"gopkg.in/macaron.v1"
	"strconv"
	"strings"
)

type ProjectForm struct {
	Id      int
	Name    string `binding:"Required"`
	Code    string `binding:"Required"`
	Remark  string
	HostIds string
}

func Index(ctx *macaron.Context) string {
	var query = models.CommonMap{}
	query["Name"] = ctx.QueryTrim("name")
	base.ParsePageAndPageSize(ctx, query)

	project := models.Project{}
	total, _ := project.Total(query)
	projects, _ := project.List(query)
	return utils.JsonResp.Success("加载成功", models.CommonMap{"total": total, "projects": projects})
}

func All(_ *macaron.Context) string {
	project := models.Project{}
	projects, _ := project.All()
	return utils.JsonResp.Success("成功", models.CommonMap{"projects": projects})
}

func Store(ctx *macaron.Context, form ProjectForm) string {
	project := models.Project{}
	project.Name = form.Name
	project.Code = form.Code
	project.Remark = form.Remark
	project.UpdateUserId = ctx.Data["uid"].(int)
	var err error
	if form.Id > 0 {
		project.Id = form.Id
		_, err = project.Update()
	} else {
		project.CreateUserId = ctx.Data["uid"].(int)
		_, err = project.Create()
	}

	if err != nil {
		return utils.JsonResp.CommonFailure("操作失败:"+err.Error(), err)
	}

	ph := models.ProjectHost{}
	_, _ = ph.RemoveForProject(project)
	for _, idStr := range strings.Split(form.HostIds, ",") {
		hostId, _ := strconv.Atoi(idStr)
		if hostId > 0 {
			ph = models.ProjectHost{HostId: hostId, ProjectId: project.Id}
			ph.Create()
		}
	}

	resetCronTask(project.Id)

	return utils.JsonResp.Success("操作成功", nil)
}

func resetCronTask(projectId int) {
	t := models.Task{}
	tasks, _ := t.ActiveListByProjectId(projectId)
	service.ServiceTask.BatchAdd(tasks)
}
