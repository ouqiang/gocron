package delaytask

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
)

type DelayForm struct {
    Protocol models.TaskType
    Type models.TaskType
    Host string
    Delay int
    Command string
    Timeout int
}

// 创建延时任务
func Create(ctx *macaron.Context)  {

}