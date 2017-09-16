package host

import (
	"fmt"
	"github.com/Unknwon/paginater"
	"github.com/go-macaron/binding"
	"github.com/ouqiang/gocron/models"
	"github.com/ouqiang/gocron/modules/logger"
	"github.com/ouqiang/gocron/modules/rpc/client"
	"github.com/ouqiang/gocron/modules/rpc/grpcpool"
	"github.com/ouqiang/gocron/modules/rpc/proto"
	"github.com/ouqiang/gocron/modules/utils"
	"github.com/ouqiang/gocron/routers/base"
	"github.com/ouqiang/gocron/service"
	"gopkg.in/macaron.v1"
	"html/template"
	"strconv"
	"strings"
)

func Index(ctx *macaron.Context) {
	hostModel := new(models.Host)
	queryParams := parseQueryParams(ctx)
	total, err := hostModel.Total(queryParams)
	hosts, err := hostModel.List(queryParams)
	if err != nil {
		logger.Error(err)
	}
	name, ok := queryParams["name"].(string)
	var safeNameHTML = ""
	if ok {
		safeNameHTML = template.HTMLEscapeString(name)
	}
	PageParams := fmt.Sprintf("id=%d&name=%s&page_size=%d",
		queryParams["Id"], safeNameHTML, queryParams["PageSize"])
	queryParams["PageParams"] = template.URL(PageParams)
	p := paginater.New(int(total), queryParams["PageSize"].(int), queryParams["Page"].(int), 5)
	ctx.Data["Pagination"] = p
	ctx.Data["Title"] = "主机列表"
	ctx.Data["Hosts"] = hosts
	ctx.Data["Params"] = queryParams
	ctx.HTML(200, "host/index")
}

func Create(ctx *macaron.Context) {
	ctx.Data["Title"] = "添加主机"
	ctx.HTML(200, "host/host_form")
}

func Edit(ctx *macaron.Context) {
	ctx.Data["Title"] = "编辑主机"
	hostModel := new(models.Host)
	id := ctx.ParamsInt(":id")
	err := hostModel.Find(id)
	if err != nil {
		logger.Errorf("获取主机详情失败#主机id-%d", id)
	}
	ctx.Data["Host"] = hostModel
	ctx.HTML(200, "host/host_form")
}

type HostForm struct {
	Id     int16
	Name   string `binding:"Required;MaxSize(64)"`
	Alias  string `binding:"Required;MaxSize(32)"`
	Port   int    `binding:"Required;Range(1-65535)"`
	Remark string
}

func (f HostForm) Error(ctx *macaron.Context, errs binding.Errors) {
	if len(errs) == 0 {
		return
	}
	json := utils.JsonResponse{}
	content := json.CommonFailure("表单验证失败, 请检测输入")

	ctx.Resp.Write([]byte(content))
}

func Store(ctx *macaron.Context, form HostForm) string {
	json := utils.JsonResponse{}
	hostModel := new(models.Host)
	id := form.Id
	nameExist, err := hostModel.NameExists(form.Name, form.Id)
	if err != nil {
		return json.CommonFailure("操作失败", err)
	}
	if nameExist {
		return json.CommonFailure("主机名已存在")
	}

	hostModel.Name = strings.TrimSpace(form.Name)
	hostModel.Alias = strings.TrimSpace(form.Alias)
	hostModel.Port = form.Port
	hostModel.Remark = strings.TrimSpace(form.Remark)
	isCreate := false
	oldHostModel := new(models.Host)
	err = oldHostModel.Find(int(id))
	if err != nil {
		return json.CommonFailure("主机不存在")
	}

	if id > 0 {
		_, err = hostModel.UpdateBean(id)
	} else {
		isCreate = true
		id, err = hostModel.Create()
	}
	if err != nil {
		return json.CommonFailure("保存失败", err)
	}

	if !isCreate {
		oldAddr := fmt.Sprintf("%s:%d", oldHostModel.Name, oldHostModel.Port)
		newAddr := fmt.Sprintf("%s:%d", hostModel.Name, hostModel.Port)
		if oldAddr != newAddr {
			grpcpool.Pool.Release(oldAddr)
		}

		taskModel := new(models.Task)
		tasks, err := taskModel.ActiveListByHostId(id)
		if err != nil {
			return json.CommonFailure("刷新任务主机信息失败", err)
		}
		serviceTask := new(service.Task)
		serviceTask.BatchAdd(tasks)
	}

	return json.Success("保存成功", nil)
}

func Remove(ctx *macaron.Context) string {
	id, err := strconv.Atoi(ctx.Params(":id"))
	json := utils.JsonResponse{}
	if err != nil {
		return json.CommonFailure("参数错误", err)
	}
	taskHostModel := new(models.TaskHost)
	exist, err := taskHostModel.HostIdExist(int16(id))
	if err != nil {
		return json.CommonFailure("操作失败", err)
	}
	if exist {
		return json.CommonFailure("有任务引用此主机，不能删除")
	}

	hostModel := new(models.Host)
	err = hostModel.Find(int(id))
	if err != nil {
		return json.CommonFailure("主机不存在")
	}

	_, err = hostModel.Delete(id)
	if err != nil {
		return json.CommonFailure("操作失败", err)
	}

	addr := fmt.Sprintf("%s:%d", hostModel.Name, hostModel.Port)
	grpcpool.Pool.Release(addr)

	return json.Success("操作成功", nil)
}

func Ping(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	hostModel := new(models.Host)
	err := hostModel.Find(id)
	json := utils.JsonResponse{}
	if err != nil || hostModel.Id <= 0 {
		return json.CommonFailure("主机不存在", err)
	}

	taskReq := &rpc.TaskRequest{}
	taskReq.Command = "echo hello"
	taskReq.Timeout = 10
	output, err := client.Exec(hostModel.Name, hostModel.Port, taskReq)
	if err != nil {
		return json.CommonFailure("连接失败-"+err.Error()+" "+output, err)
	}

	return json.Success("连接成功", nil)
}

// 解析查询参数
func parseQueryParams(ctx *macaron.Context) models.CommonMap {
	var params models.CommonMap = models.CommonMap{}
	params["Id"] = ctx.QueryInt("id")
	params["Name"] = ctx.QueryTrim("name")
	base.ParsePageAndPageSize(ctx, params)

	return params
}
