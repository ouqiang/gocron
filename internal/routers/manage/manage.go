package manage

import (
	"encoding/json"

	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"gopkg.in/macaron.v1"
)

// region slack

func EditSlack(ctx *macaron.Context) {
	ctx.Data["Title"] = "Slack配置"
	settingModel := new(models.Setting)
	slack, err := settingModel.Slack()
	if err != nil {
		logger.Error(err)
	}
	ctx.Data["Slack"] = slack
	ctx.HTML(200, "manage/slack")
}

func Slack(ctx *macaron.Context) string {
	settingModel := new(models.Setting)
	slack, err := settingModel.Slack()
	if err != nil {
		logger.Error(err)
	}
	json := utils.JsonResponse{}

	return json.Success("", slack)
}

func UpdateSlackUrl(ctx *macaron.Context) string {
	url := ctx.QueryTrim("url")
	settingModel := new(models.Setting)
	_, err := settingModel.UpdateSlackUrl(url)

	return utils.JsonResponseByErr(err)
}

func CreateSlackChannel(ctx *macaron.Context) string {
	channel := ctx.QueryTrim("channel")
	settingModel := new(models.Setting)
	if settingModel.IsChannelExist(channel) {
		json := utils.JsonResponse{}

		return json.CommonFailure("Channel已存在")
	}
	_, err := settingModel.CreateChannel(channel)

	return utils.JsonResponseByErr(err)
}

func RemoveSlackChannel(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	settingModel := new(models.Setting)
	_, err := settingModel.RemoveChannel(id)

	return utils.JsonResponseByErr(err)
}

// endregion

// region 邮件

func EditMail(ctx *macaron.Context) {
	ctx.Data["Title"] = "邮件配置"
	settingModel := new(models.Setting)
	mail, err := settingModel.Mail()
	if err != nil {
		logger.Error(err)
	}
	ctx.Data["Mail"] = mail
	ctx.HTML(200, "manage/mail")
}

func Mail(ctx *macaron.Context) string {
	settingModel := new(models.Setting)
	mail, err := settingModel.Mail()
	if err != nil {
		logger.Error(err)
	}

	json := utils.JsonResponse{}

	return json.Success("", mail)
}

type MailServerForm struct {
	Host     string `binding:"Required;MaxSize(100)"`
	Port     int    `binding:"Required;Range(1-65535)"`
	User     string `binding:"Required;MaxSize(64);Email"`
	Password string `binding:"Required;MaxSize(64)"`
}

func UpdateMailServer(ctx *macaron.Context, form MailServerForm) string {
	jsonByte, _ := json.Marshal(form)
	settingModel := new(models.Setting)
	_, err := settingModel.UpdateMailServer(string(jsonByte))

	return utils.JsonResponseByErr(err)
}

func ClearMailServer(ctx *macaron.Context) string {
	jsonByte, _ := json.Marshal(MailServerForm{})
	settingModel := new(models.Setting)
	_, err := settingModel.UpdateMailServer(string(jsonByte))

	return utils.JsonResponseByErr(err)
}

func CreateMailUser(ctx *macaron.Context) string {
	username := ctx.QueryTrim("username")
	email := ctx.QueryTrim("email")
	settingModel := new(models.Setting)
	if username == "" || email == "" {
		json := utils.JsonResponse{}

		return json.CommonFailure("用户名、邮箱均不能为空")
	}
	_, err := settingModel.CreateMailUser(username, email)

	return utils.JsonResponseByErr(err)
}

func RemoveMailUser(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	settingModel := new(models.Setting)
	_, err := settingModel.RemoveMailUser(id)

	return utils.JsonResponseByErr(err)
}

// endregion
