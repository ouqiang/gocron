package notify

import (
    "time"
)

var SlackUrl string

type Message map[string]interface{}
type Notifiable interface {
    send(msg Message)
}

var queue chan Message = make(chan Message, 100)

func init()  {
    go run()
}

// 把消息推入队列
func Push(msg Message) {
    queue <- msg
}

func run() {
    slack := new(Slack)
    for msg := range queue {
        // 根据任务配置执行相应通知
        go slack.Send(msg)
        time.Sleep(1 * time.Second)
    }
}