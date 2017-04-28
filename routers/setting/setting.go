package setting

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/ouqiang/gocron/models"
    "github.com/ouqiang/gocron/modules/logger"
)

func EditSlack(ctx *macaron.Context)  {
    ctx.Data["Title"] = "slack配置"
    settingModel := new(models.Setting)
    url, err := settingModel.SlackUrl()
    if err != nil {
        logger.Error(err)
    }
    ctx.Data["SlackUrl"] = url
    ctx.HTML(200, "setting/slack")
}

func StoreSlack(ctx *macaron.Context) string {
    url := ctx.QueryTrim("url")
    settingModel := new(models.Setting)
    _, err := settingModel.UpdateSlackUrl(url)
    json := utils.JsonResponse{}
    if err != nil {
        return json.CommonFailure(utils.FailureContent, err)
    }

    return json.Success(utils.SuccessContent, nil)
}