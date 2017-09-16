package models

import (
	"encoding/json"
)

type Setting struct {
	Id    int    `xorm:"int pk autoincr"`
	Code  string `xorm:"varchar(32) notnull"`
	Key   string `xorm:"varchar(64) notnull"`
	Value string `xorm:"varchar(4096) notnull default '' "`
}

const SlackCode = "slack"
const SlackUrlKey = "url"
const SlackChannelKey = "channel"

const MailCode = "mail"
const MailServerKey = "server"
const MailUserKey = "user"

// 初始化基本字段 邮件、slack等
func (setting *Setting) InitBasicField() {
	setting.Code = SlackCode
	setting.Key = SlackUrlKey
	Db.Insert(setting)

	setting.Id = 0
	setting.Code = MailCode
	setting.Key = MailServerKey
	Db.Insert(setting)
}

// region slack配置

type Slack struct {
	Url      string
	Channels []Channel
}

type Channel struct {
	Id   int
	Name string
}

func (setting *Setting) Slack() (Slack, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", SlackCode).Find(&list)
	slack := Slack{Url: "", Channels: make([]Channel, 0)}
	if err != nil {
		return slack, err
	}

	setting.formatSlack(list, &slack)

	return slack, err
}

func (setting *Setting) formatSlack(list []Setting, slack *Slack) {
	for _, v := range list {
		if v.Key == SlackUrlKey {
			slack.Url = v.Value
			continue
		}

		slack.Channels = append(slack.Channels, Channel{
			v.Id, v.Value,
		})
	}
}

// 更新slack webhook url
func (setting *Setting) UpdateSlackUrl(url string) (int64, error) {
	setting.Value = url

	return Db.Cols("value").Update(setting, Setting{Code: SlackCode, Key: SlackUrlKey})
}

// 创建slack渠道
func (setting *Setting) CreateChannel(channel string) (int64, error) {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel

	return Db.Insert(setting)
}

func (setting *Setting) IsChannelExist(channel string) bool {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel

	count, _ := Db.Count(setting)

	return count > 0
}

// 删除slack渠道
func (setting *Setting) RemoveChannel(id int) (int64, error) {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Id = id
	return Db.Delete(setting)
}

// endregion

type Mail struct {
	Host      string
	Port      int
	User      string
	Password  string
	MailUsers []MailUser
}

type MailUser struct {
	Id       int
	Username string
	Email    string
}

// region 邮件配置
func (setting *Setting) Mail() (Mail, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", MailCode).Find(&list)
	mail := Mail{MailUsers: make([]MailUser, 0)}
	if err != nil {
		return mail, err
	}

	setting.formatMail(list, &mail)

	return mail, err
}

func (setting *Setting) formatMail(list []Setting, mail *Mail) {
	mailUser := MailUser{}
	for _, v := range list {
		if v.Key == MailServerKey {
			json.Unmarshal([]byte(v.Value), mail)
			continue
		}
		json.Unmarshal([]byte(v.Value), &mailUser)
		mailUser.Id = v.Id
		mail.MailUsers = append(mail.MailUsers, mailUser)
	}
}

func (setting *Setting) UpdateMailServer(config string) (int64, error) {
	setting.Value = config
	return Db.Cols("value").Update(setting, Setting{Code: MailCode, Key: MailServerKey})
}

func (setting *Setting) CreateMailUser(username, email string) (int64, error) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	mailUser := MailUser{0, username, email}
	jsonByte, err := json.Marshal(mailUser)
	if err != nil {
		return 0, err
	}
	setting.Value = string(jsonByte)

	return Db.Insert(setting)
}

func (setting *Setting) RemoveMailUser(id int) (int64, error) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	setting.Id = id
	return Db.Delete(setting)
}

// endregion
