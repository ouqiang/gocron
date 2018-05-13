package notify

// 发送消息到slack

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/httpclient"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
)

type Slack struct{}

func (slack *Slack) Send(msg Message) {
	model := new(models.Setting)
	slackSetting, err := model.Slack()
	if err != nil {
		logger.Error("#slack#从数据库获取slack配置失败", err)
		return
	}
	if slackSetting.Url == "" {
		logger.Error("#slack#webhook-url为空")
		return
	}
	if len(slackSetting.Channels) == 0 {
		logger.Error("#slack#channels配置为空")
		return
	}
	logger.Debugf("%+v", slackSetting)
	channels := slack.getActiveSlackChannels(slackSetting, msg)
	logger.Debugf("%+v", channels)
	msg["content"] = parseNotifyTemplate(slackSetting.Template, msg)
	msg["content"] = html.UnescapeString(msg["content"].(string))
	for _, channel := range channels {
		slack.send(msg, slackSetting.Url, channel)
	}
}

func (slack *Slack) send(msg Message, slackUrl string, channel string) {
	formatBody := slack.format(msg["content"].(string), channel)
	timeout := 30
	maxTimes := 3
	i := 0
	for i < maxTimes {
		resp := httpclient.PostJson(slackUrl, formatBody, timeout)
		if resp.StatusCode == 200 {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			logger.Errorf("slack#发送消息失败#%s#消息内容-%s", resp.Body, msg["content"])
		}
	}
}

func (slack *Slack) getActiveSlackChannels(slackSetting models.Slack, msg Message) []string {
	taskReceiverIds := strings.Split(msg["task_receiver_id"].(string), ",")
	channels := []string{}
	for _, v := range slackSetting.Channels {
		if utils.InStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
			channels = append(channels, v.Name)
		}
	}

	return channels
}

// 格式化消息内容
func (slack *Slack) format(content string, channel string) string {
	content = utils.EscapeJson(content)
	specialChars := []string{"&", "<", ">"}
	replaceChars := []string{"&amp;", "&lt;", "&gt;"}
	content = utils.ReplaceStrings(content, specialChars, replaceChars)

	return fmt.Sprintf(`{"text":"%s","username":"gocron", "channel":"%s"}`, content, channel)
}
