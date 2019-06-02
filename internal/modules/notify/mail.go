package notify

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/utils"
)

// @author qiang.ou<qingqianludao@gmail.com>
// @date 2017/5/1-00:19

type Mail struct {
}

func (mail *Mail) Send(msg Message) {
	model := new(models.Setting)
	mailSetting, err := model.Mail()
	logger.Debugf("%+v", mailSetting)
	if err != nil {
		logger.Error("#mail#从数据库获取mail配置失败", err)
		return
	}
	if mailSetting.Host == "" {
		logger.Error("#mail#Host为空")
		return
	}
	if mailSetting.Port == 0 {
		logger.Error("#mail#Port为空")
		return
	}
	if mailSetting.User == "" {
		logger.Error("#mail#User为空")
		return
	}
	if mailSetting.Password == "" {
		logger.Error("#mail#Password为空")
		return
	}
	msg["content"] = parseNotifyTemplate(mailSetting.Template, msg)
	toUsers := mail.getActiveMailUsers(mailSetting, msg)
	mail.send(mailSetting, toUsers, msg)
}

func (mail *Mail) send(mailSetting models.Mail, toUsers []string, msg Message) {
	body := msg["content"].(string)
	body = strings.Replace(body, "\n", "<br>", -1)
	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("From", mailSetting.User)
	gomailMessage.SetHeader("To", toUsers...)
	gomailMessage.SetHeader("Subject", "gocron-定时任务通知")
	gomailMessage.SetBody("text/html", body)
	mailer := gomail.NewDialer(mailSetting.Host, mailSetting.Port,
		mailSetting.User, mailSetting.Password)
	maxTimes := 3
	i := 0
	for i < maxTimes {
		err := mailer.DialAndSend(gomailMessage)
		if err == nil {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			logger.Errorf("mail#发送消息失败#%s#消息内容-%s", err.Error(), msg["content"])
		}
	}
}

func (mail *Mail) getActiveMailUsers(mailSetting models.Mail, msg Message) []string {
	taskReceiverIds := strings.Split(msg["task_receiver_id"].(string), ",")
	users := []string{}
	for _, v := range mailSetting.MailUsers {
		if utils.InStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
			users = append(users, v.Email)
		}
	}

	return users
}
