package models

import (
    "time"
    "github.com/go-xorm/xorm"
)


// 延迟任务
type DelayTask struct {
    Id       int64     `xorm:"bigint pk autoincr"`
    Url      string    `xorm:"varchar(128) not null"`
    Params   string    `xorm:"varchar(256) not null default '' "`
    Delay  int       `xorm:"mediumint notnull default 0"` // 延迟时间
    Status   Status    `xorm:"tinyint notnull index(u_status_created) default 5"` // 状态 0:执行失败 1:执行中 2:执行成功 5: 待执行
    Created  time.Time `xorm:"datetime notnull created index(u_status_created)"`
    Updated  time.Time `xorm:"datetime updated"`
    BaseModel `xorm:"-"`
}

func (task *DelayTask) Create() (insertId int64, err error) {
    _, err = Db.Insert(task)
    if err == nil {
        insertId = task.Id
    }

    return
}

// 获取所有待执行任务
func (task *DelayTask) ActiveList(endTime time.Time) ([]DelayTask, error) {
    list := make([]DelayTask, 0)
    fields := "id,url,params,delay,created"
    err := Db.Where("status = ? AND created <= ?", Waiting, endTime.Format(DefaultTimeFormat)).Cols(fields).Limit(task.PageSize, task.pageLimitOffset()).Find(&list)

    return list, err
}

// 获取待执行任务数量
func (task *DelayTask) ActiveNum(endTime time.Time) (int, error) {
    count ,err := Db.Where("status = ? AND created <= ?", Waiting, endTime.Format(DefaultTimeFormat)).Count(task)

    return int(count), err
}

func (task *DelayTask) List(params CommonMap) ([]DelayTask, error) {
    task.parsePageAndPageSize(params)
    list := make([]DelayTask, 0)
    session := Db.Desc("id")
    task.parseWhere(session, params)
    err := session.Limit(task.PageSize, task.pageLimitOffset()).Find(&list)

    return list, err
}


// 更新任务状态
func (task *DelayTask) UpdateStatus(id int64, status Status) (int64, error) {
    return Db.Table(task).Id(id).Update(CommonMap{
        "status": status,
    })
}

// 解析where
func (task *DelayTask) parseWhere(session *xorm.Session, params CommonMap)  {
    if len(params) == 0 {
        return
    }
    status, ok := params["Status"]
    if ok && status.(int) > -1 {
        session.And("status = ?", status)
    }
}

// 删除N个月前的日志
func (task *DelayTask) Remove(id int) (int64, error) {
    t := time.Now().AddDate(0, -id, 0)
    return Db.Where("created <= ?", t.Format(DefaultTimeFormat)).Delete(task)
}

func (task *DelayTask) Total(params CommonMap) (int64, error) {
    session := Db.NewSession()
    task.parseWhere(session, params)
    return session.Count(task)
}
