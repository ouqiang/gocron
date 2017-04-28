package models

type Setting struct  {
    Id        int16     `xorm:"smallint pk autoincr"`
    Code      string    `xorm:"varchar(32) notnull"`
    Key       string    `xorm:"varchar(64) notnull"`
    Value     string    `xorm:"varchar(4096) notnull default '' "`
}

const SlackCode = "slack"
const SlackKey  = "url"

// 初始化基本字段 邮件、slack等
func (setting *Setting) InitBasicField() (error) {
    setting.Code = "slack";
    setting.Key = "url"
    setting.Value = ""
    _, err := Db.Insert(setting)

    return err
}

func (setting *Setting) SlackUrl() (string, error)  {
    setting.slackCondition()
    _, err := Db.Get(setting)

    return setting.Value, err
}

func (setting *Setting) UpdateSlackUrl(url string) (int64, error) {
    setting.slackCondition()
    setting.Value = url

    return setting.UpdateBean()
}

func (setting *Setting) slackCondition()  {
    setting.Code = SlackCode
    setting.Key = SlackKey
}

func (setting *Setting) UpdateBean() (int64, error)  {
    return Db.Cols("code,key,value").Update(setting)
}