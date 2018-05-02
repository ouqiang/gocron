package models

import (
	"time"
)

// 用户登录日志

type LoginLog struct {
	Id        int       `json:"id" xorm:"pk autoincr notnull "`
	Username  string    `json:"username" xorm:"varchar(32) notnull"`
	Ip        string    `json:"ip" xorm:"varchar(15) not null"`
	Created   time.Time `json:"created" xorm:"datetime notnull created"`
	BaseModel `json:"-" xorm:"-"`
}

func (log *LoginLog) Create() (insertId int, err error) {
	_, err = Db.Insert(log)
	if err == nil {
		insertId = log.Id
	}

	return
}

func (log *LoginLog) List(params CommonMap) ([]LoginLog, error) {
	log.parsePageAndPageSize(params)
	list := make([]LoginLog, 0)
	err := Db.Desc("id").Limit(log.PageSize, log.pageLimitOffset()).Find(&list)

	return list, err
}

func (log *LoginLog) Total() (int64, error) {
	return Db.Count(log)
}
