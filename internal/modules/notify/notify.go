package notify

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/ouqiang/gocron/internal/modules/logger"
)

type Message map[string]interface{}

type Notifiable interface {
	Send(msg Message)
}

var queue = make(chan Message, 100)

func init() {
	go run()
}

// 把消息推入队列
func Push(msg Message) {
	queue <- msg
}

func run() {
	for msg := range queue {
		// 根据任务配置发送通知
		taskType, taskTypeOk := msg["task_type"]
		_, taskReceiverIdOk := msg["task_receiver_id"]
		_, nameOk := msg["name"]
		_, outputOk := msg["output"]
		_, statusOk := msg["status"]
		if !taskTypeOk || !taskReceiverIdOk || !nameOk || !outputOk || !statusOk {
			logger.Errorf("#notify#参数不完整#%+v", msg)
			continue
		}
		msg["content"] = fmt.Sprintf("============\n============\n============\n任务名称: %s\n状态: %s\n输出:\n %s\n", msg["name"], msg["status"], msg["output"])
		logger.Debugf("%+v", msg)
		switch taskType.(int8) {
		case 1:
			// 邮件
			mail := Mail{}
			go mail.Send(msg)
		case 2:
			// Slack
			slack := Slack{}
			go slack.Send(msg)
		case 3:
			// WebHook
			webHook := WebHook{}
			go webHook.Send(msg)
		}
		time.Sleep(1 * time.Second)
	}
}

func parseNotifyTemplate(notifyTemplate string, msg Message) string {
	tmpl, err := template.New("notify").Parse(notifyTemplate)
	if err != nil {
		return fmt.Sprintf("解析通知模板失败: %s", err)
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, map[string]interface{}{
		"TaskId":   msg["task_id"],
		"TaskName": msg["name"],
		"Status":   msg["status"],
		"Result":   msg["output"],
	})

	return buf.String()
}
