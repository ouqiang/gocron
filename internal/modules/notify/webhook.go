package notify

import (
	"html"
	"time"

	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/httpclient"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
)

type WebHook struct{}

func (webHook *WebHook) Send(msg Message) {
	model := new(models.Setting)
	webHookSetting, err := model.Webhook()
	if err != nil {
		logger.Error("#webHook#从数据库获取webHook配置失败", err)
		return
	}
	if webHookSetting.Url == "" {
		logger.Error("#webHook#webhook-url为空")
		return
	}
	logger.Debugf("%+v", webHookSetting)
	msg["name"] = utils.EscapeJson(msg["name"].(string))
	msg["output"] = utils.EscapeJson(msg["output"].(string))
	msg["content"] = parseNotifyTemplate(webHookSetting.Template, msg)
	msg["content"] = html.UnescapeString(msg["content"].(string))
	webHook.send(msg, webHookSetting.Url)
}

func (webHook *WebHook) send(msg Message, url string) {
	content := msg["content"].(string)
	timeout := 30
	maxTimes := 3
	i := 0
	for i < maxTimes {
		resp := httpclient.PostJson(url, content, timeout)
		if resp.StatusCode == 200 {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			logger.Errorf("webHook#发送消息失败#%s#消息内容-%s", resp.Body, msg["content"])
		}
	}
}
