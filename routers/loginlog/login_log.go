package loginlog

import (
	"fmt"
	"html/template"

	"github.com/Unknwon/paginater"
	"github.com/ouqiang/gocron/models"
	"github.com/ouqiang/gocron/modules/logger"
	"github.com/ouqiang/gocron/routers/base"
	"gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) {
	loginLogModel := new(models.LoginLog)
	params := models.CommonMap{}
	base.ParsePageAndPageSize(ctx, params)
	total, err := loginLogModel.Total()
	loginLogs, err := loginLogModel.List(params)
	if err != nil {
		logger.Error(err)
	}
	PageParams := fmt.Sprintf("page_size=%d", params["PageSize"])
	params["PageParams"] = template.URL(PageParams)
	p := paginater.New(int(total), params["PageSize"].(int), params["Page"].(int), 5)
	ctx.Data["Pagination"] = p
	ctx.Data["Title"] = "登录日志"
	ctx.Data["LoginLogs"] = loginLogs
	ctx.Data["Params"] = params
	ctx.HTML(200, "manage/login_log")
}
