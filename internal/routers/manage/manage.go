package manage

import (
	"encoding/json"

	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"gopkg.in/macaron.v1"
)

func Slack(ctx *macaron.Context) string {
	settingModel := new(models.Setting)
	slack, err := settingModel.Slack()
	jsonResp := utils.JsonResponse{}
	if err != nil {
		logger.Error(err)
		return jsonResp.Success(utils.SuccessContent, nil)

	}

	return jsonResp.Success(utils.SuccessContent, slack)
}

func UpdateSlack(ctx *macaron.Context) string {
	url := ctx.QueryTrim("url")
	template := ctx.QueryTrim("template")
	settingModel := new(models.Setting)
	err := settingModel.UpdateSlack(url, template)

	return utils.JsonResponseByErr(err)
}

func CreateSlackChannel(ctx *macaron.Context) string {
	channel := ctx.QueryTrim("channel")
	settingModel := new(models.Setting)
	if settingModel.IsChannelExist(channel) {
		jsonResp := utils.JsonResponse{}

		return jsonResp.CommonFailure("Channel已存在")
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
func Mail(ctx *macaron.Context) string {
	settingModel := new(models.Setting)
	mail, err := settingModel.Mail()
	jsonResp := utils.JsonResponse{}
	if err != nil {
		logger.Error(err)
		return jsonResp.Success(utils.SuccessContent, nil)
	}

	return jsonResp.Success("", mail)
}

type MailServerForm struct {
	Host     string `binding:"Required;MaxSize(100)"`
	Port     int    `binding:"Required;Range(1-65535)"`
	User     string `binding:"Required;MaxSize(64);Email"`
	Password string `binding:"Required;MaxSize(64)"`
}

func UpdateMail(ctx *macaron.Context, form MailServerForm) string {
	jsonByte, _ := json.Marshal(form)
	settingModel := new(models.Setting)

	template := ctx.QueryTrim("template")
	err := settingModel.UpdateMail(string(jsonByte), template)

	return utils.JsonResponseByErr(err)
}

func CreateMailUser(ctx *macaron.Context) string {
	username := ctx.QueryTrim("username")
	email := ctx.QueryTrim("email")
	settingModel := new(models.Setting)
	if username == "" || email == "" {
		jsonResp := utils.JsonResponse{}

		return jsonResp.CommonFailure("用户名、邮箱均不能为空")
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

func WebHook(ctx *macaron.Context) string {
	settingModel := new(models.Setting)
	webHook, err := settingModel.Webhook()
	jsonResp := utils.JsonResponse{}
	if err != nil {
		logger.Error(err)
		return jsonResp.Success(utils.SuccessContent, nil)
	}

	return jsonResp.Success("", webHook)
}

func UpdateWebHook(ctx *macaron.Context) string {
	url := ctx.QueryTrim("url")
	template := ctx.QueryTrim("template")
	settingModel := new(models.Setting)
	err := settingModel.UpdateWebHook(url, template)

	return utils.JsonResponseByErr(err)
}

// endregion
