package notify
// 发送消息到slack

import (
    "fmt"
    "github.com/ouqiang/gocron/modules/httpclient"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/ouqiang/gocron/modules/utils"
)

type Slack struct {}

func (slack *Slack) Send(msg Message)  {
    name, nameOk := msg["name"]
    statusName, statusOk := msg["status"]
    content, contentOk := msg["output"]
    if SlackUrl == "" {
        logger.Error("slack#webhooks-url为空")
        return;
    }
    if !nameOk || !statusOk || !contentOk {
        logger.Error("slack#消息字段不存在")
        return
    }
    body := fmt.Sprintf("============\n============\n============\n任务名称: %s\n状态: %s\n输出:\n %s\n", name, statusName, content)
    formatBody := slack.format(body)
    timeout := 30
    maxTimes := 3
    i := 0
    for i < maxTimes {
        resp := httpclient.PostBody(SlackUrl, formatBody, timeout)
        if resp.StatusCode == 200 {
            break;
        }
        i += 1
        if i < maxTimes {
            logger.Error("slack#发送消息失败#%s#消息内容-%s", resp.Body, body)
        }
    }
}

// 格式化消息内容
func (slack *Slack) format(content string) string  {
    content = utils.EscapeJson(content)
    specialChars := []string{"&", "<", ">"}
    replaceChars := []string{"&amp;", "&lt;", "&gt;"}
    content = utils.ReplaceStrings(content, specialChars, replaceChars)

    return fmt.Sprintf(`{"text":"%s","username":"监控"}`, content)
}