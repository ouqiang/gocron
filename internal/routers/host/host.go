package host

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-macaron/binding"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/rpc/client"
	"github.com/ouqiang/gocron/internal/modules/rpc/grpcpool"
	"github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/gocron/internal/routers/base"
	"github.com/ouqiang/gocron/internal/service"
	macaron "gopkg.in/macaron.v1"
)

const testConnectionCommand = "echo hello"
const testConnectionTimeout = 5

// Index 主机列表
func Index(ctx *macaron.Context) string {
	hostModel := new(models.Host)
	queryParams := parseQueryParams(ctx)
	total, err := hostModel.Total(queryParams)
	if err != nil {
		logger.Error(err)
	}
	hosts, err := hostModel.List(queryParams)
	if err != nil {
		logger.Error(err)
	}

	jsonResp := utils.JsonResponse{}

	return jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  hosts,
	})
}

// All 获取所有主机
func All(ctx *macaron.Context) string {
	hostModel := new(models.Host)
	hostModel.PageSize = -1
	hosts, err := hostModel.List(models.CommonMap{})
	if err != nil {
		logger.Error(err)
	}

	jsonResp := utils.JsonResponse{}

	return jsonResp.Success(utils.SuccessContent, hosts)
}

// Detail 主机详情
func Detail(ctx *macaron.Context) string {
	hostModel := new(models.Host)
	id := ctx.ParamsInt(":id")
	err := hostModel.Find(id)
	jsonResp := utils.JsonResponse{}
	if err != nil || hostModel.Id == 0 {
		logger.Errorf("获取主机详情失败#主机id-%d", id)
		return jsonResp.Success(utils.SuccessContent, nil)
	}

	return jsonResp.Success(utils.SuccessContent, hostModel)
}

type HostForm struct {
	Id     int16
	Name   string `binding:"Required;MaxSize(64)"`
	Alias  string `binding:"Required;MaxSize(32)"`
	Port   int    `binding:"Required;Range(1-65535)"`
	Remark string
}

// Error 表单验证错误处理
func (f HostForm) Error(ctx *macaron.Context, errs binding.Errors) {
	if len(errs) == 0 {
		return
	}
	json := utils.JsonResponse{}
	content := json.CommonFailure("表单验证失败, 请检测输入")
	ctx.Write([]byte(content))
}

// Store 保存、修改主机信息
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
		service.ServiceTask.BatchAdd(tasks)
	}

	return json.Success("保存成功", nil)
}

// Remove 删除主机
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

// Ping 测试主机是否可连接
func Ping(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	hostModel := new(models.Host)
	err := hostModel.Find(id)
	json := utils.JsonResponse{}
	if err != nil || hostModel.Id <= 0 {
		return json.CommonFailure("主机不存在", err)
	}

	taskReq := &rpc.TaskRequest{}
	taskReq.Command = testConnectionCommand
	taskReq.Timeout = testConnectionTimeout
	output, err := client.Exec(hostModel.Name, hostModel.Port, taskReq)
	if err != nil {
		return json.CommonFailure("连接失败-"+err.Error()+" "+output, err)
	}

	return json.Success("连接成功", nil)
}

// 解析查询参数
func parseQueryParams(ctx *macaron.Context) models.CommonMap {
	var params = models.CommonMap{}
	params["Id"] = ctx.QueryInt("id")
	params["Name"] = ctx.QueryTrim("name")
	base.ParsePageAndPageSize(ctx, params)

	return params
}
