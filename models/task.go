package models

import (
    "time"
    "github.com/go-xorm/xorm"
    "errors"
)

type TaskProtocol int8

const (
    TaskHTTP TaskProtocol = iota + 1 // HTTP协议
    TaskRPC  // RPC方式执行命令
)

// 任务
type Task struct {
    Id       int       `xorm:"int pk autoincr"`
    Name     string    `xorm:"varchar(32) notnull"`              // 任务名称
    Spec     string    `xorm:"varchar(64) notnull"`              // crontab
    Protocol TaskProtocol  `xorm:"tinyint notnull"`              // 协议 1:http 2:系统命令
    Command  string    `xorm:"varchar(256) notnull"`             // URL地址或shell命令
    Timeout  int       `xorm:"mediumint notnull default 0"`      // 任务执行超时时间(单位秒),0不限制
    Multi    int8      `xorm:"tinyint notnull default 1"`        // 是否允许多实例运行
    RetryTimes int8    `xorm:"tinyint notnull default 0"`         // 重试次数
    HostId   int16    `xorm:"smallint notnull default 0"`        // RPC host id，
    NotifyStatus int8  `xorm:"smallint notnull default 1"`       // 任务执行结束是否通知 0: 不通知 1: 失败通知 2: 执行结束通知
    NotifyType int8 `xorm:"smallint notnull default 0"`  // 通知类型 1: 邮件 2: slack
    NotifyReceiverId string `xorm:"varchar(256) notnull default '' "` // 通知接受者ID, setting表主键ID，多个ID逗号分隔
    Remark   string    `xorm:"varchar(100) notnull default ''"`  // 备注
    Status   Status    `xorm:"tinyint notnull default 0"`        // 状态 1:正常 0:停止
    Created  time.Time `xorm:"datetime notnull created"`         // 创建时间
    Deleted  time.Time `xorm:"datetime deleted"`                 // 删除时间
    BaseModel `xorm:"-"`
}

type TaskHost struct {
    Task `xorm:"extends"`
    Name string
    Port int
    Alias string
}

func (TaskHost) TableName() string  {
    return TablePrefix + "task"
}

// 新增
func (task *Task) Create() (insertId int, err error) {
    _, err = Db.Insert(task)
    if err == nil {
        insertId = task.Id
    }

    return
}

// 新增测试任务
func (task *Task) CreateTestTask() {
    // HTTP任务
    task.Name = "测试HTTP任务"
    task.Protocol = TaskHTTP
    task.Spec = "*/30 * * * * *"
    // 查询IP地址区域信息
    task.Command = "http://ip.taobao.com/service/getIpInfo.php?ip=117.27.140.253"
    task.Status = Enabled
    task.Create()
}

func (task *Task) UpdateBean(id int) (int64, error)  {
    return Db.ID(id).Cols("name,spec,protocol,command,timeout,multi,retry_times,host_id,remark,notify_status,notify_type,notify_receiver_id").Update(task)
}

// 更新
func (task *Task) Update(id int, data CommonMap) (int64, error) {
    return Db.Table(task).ID(id).Update(data)
}

// 删除
func (task *Task) Delete(id int) (int64, error) {
    return Db.Id(id).Delete(task)
}

// 禁用
func (task *Task) Disable(id int) (int64, error) {
    return task.Update(id, CommonMap{"status": Disabled})
}

// 激活
func (task *Task) Enable(id int) (int64, error) {
    return task.Update(id, CommonMap{"status": Enabled})
}

// 获取所有激活任务
func (task *Task) ActiveList() ([]TaskHost, error) {
    list := make([]TaskHost, 0)
    fields := "t.*, host.alias,host.name,host.port"
    err := Db.Alias("t").Join("LEFT", hostTableName(), "t.host_id=host.id").Where("t.status = ?", Enabled).Cols(fields).Find(&list)

    return list, err
}

// 获取某个主机下的所有激活任务
func (task *Task) ActiveListByHostId(hostId int16) ([]TaskHost, error) {
    list := make([]TaskHost, 0)
    fields := "t.*, host.alias,host.name,host.port"
    err := Db.Alias("t").Join("LEFT", hostTableName(), "t.host_id=host.id").Where("t.status = ? AND t.host_id = ?", Enabled, hostId).Cols(fields).Find(&list)

    return list, err
}

// 判断主机id是否有引用
func (task *Task) HostIdExist(hostId int16) (bool, error) {
    count, err := Db.Where("host_id = ?", hostId).Count(task);

    return count > 0, err
}

// 判断任务名称是否存在
func (task *Task) NameExist(name string, id int) (bool, error)  {
    if id > 0 {
        count, err := Db.Where("name = ? AND status = ? AND id != ?", name, Enabled, id).Count(task);
        return count > 0, err
    }
    count, err := Db.Where("name = ? AND status = ?", name, Enabled).Count(task);

    return count > 0, err
}

func (task *Task) GetStatus(id int) (Status, error) {
    exist, err := Db.Id(id).Get(task)
    if err != nil {
        return 0, err
    }
    if !exist {
        return 0, errors.New("not exist")
    }

    return task.Status, nil
}

func(task *Task) Detail(id int) (TaskHost, error)  {
    taskHost := TaskHost{}
    fields := "t.*, host.alias,host.name,host.port"
    _, err := Db.Alias("t").Join("LEFT", hostTableName(), "t.host_id=host.id").Where("t.id=?", id).Cols(fields).Get(&taskHost)

    return taskHost, err
}

func (task *Task) List(params CommonMap) ([]TaskHost, error) {
    task.parsePageAndPageSize(params)
    list := make([]TaskHost, 0)
    fields := "t.*, host.alias,host.name"
    session := Db.Alias("t").Join("LEFT", hostTableName(), "t.host_id=host.id")
    task.parseWhere(session, params)
    err := session.Cols(fields).Desc("t.id").Limit(task.PageSize, task.pageLimitOffset()).Find(&list)

    return list, err
}

func (task *Task) Total(params CommonMap) (int64, error) {
    session := Db.Alias("t").Join("LEFT", hostTableName(), "t.host_id=host.id")
    task.parseWhere(session, params)
    return session.Count(task)
}

// 解析where
func (task *Task) parseWhere(session *xorm.Session, params CommonMap)  {
    if len(params) == 0 {
        return
    }
    id, ok := params["Id"]
    if ok && id.(int) > 0 {
        session.And("t.id = ?", id)
    }
    hostId, ok := params["HostId"]
    if ok && hostId.(int) > 0 {
        session.And("host_id = ?", hostId)
    }
    name, ok := params["Name"]
    if ok && name.(string) != "" {
        session.And("t.name LIKE ?", "%" + name.(string) + "%")
    }
    protocol, ok := params["Protocol"]
    if ok && protocol.(int) > 0 {
        session.And("protocol = ?", protocol)
    }
    status, ok := params["Status"]
    if ok && status.(int) > -1 {
        session.And("status = ?", status)
    }
}

func hostTableName() []string {
    return []string{TablePrefix + "host", "host"}
}