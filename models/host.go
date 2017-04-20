package models

import "github.com/ouqiang/gocron/modules/ssh"



// 主机
type Host struct {
    Id        int16     `xorm:"smallint pk autoincr"`
    Name      string    `xorm:"varchar(128) notnull"`             // 主机名称
    Alias     string    `xorm:"varchar(32) notnull default '' "`  // 主机别名
    Username  string    `xorm:"varchar(32) notnull default '' "`  // ssh 用户名
    Password  string    `xorm:"varchar(64) notnull default ''"`   // ssh 密码
    Port      int       `xorm:"notnull default 22"`               // 主机端口
    Remark    string    `xorm:"varchar(512) notnull default '' "` // 备注
    AuthType  ssh.HostAuthType      `xorm:"tinyint notnull default 1"` // 认证方式 1: 密码 2: 公钥
    PrivateKey string   `xorm:"varchar(4096) notnull default '' "` // 私钥
    Page      int       `xorm:"-"`
    PageSize  int       `xorm:"-"`
}

// 新增
func (host *Host) Create() (insertId int16, err error) {
    _, err = Db.Insert(host)
    if err == nil {
        insertId = host.Id
    }

    return
}

func (host *Host) UpdateBean() (int64, error)  {
    return Db.Cols("name,alias,username,password,port,remark,auth_type,private_key").Update(host)
}


// 更新
func (host *Host) Update(id int, data CommonMap) (int64, error) {
    return Db.Table(host).ID(id).Update(data)
}

// 删除
func (host *Host) Delete(id int) (int64, error) {
    return Db.Id(id).Delete(host)
}

func (host *Host) Find(id int) error {
    _, err := Db.Id(id).Get(host)

    return err
}

func (host *Host) NameExists(name string, id int16) (bool, error)  {
    if id == 0 {
        count, err := Db.Where("name = ?", name).Count(host);
        return count > 0, err
    }

    count, err := Db.Where("name = ? AND id != ?", name, id).Count(host);
    return count > 0, err
}

func (host *Host) List() ([]Host, error) {
    host.parsePageAndPageSize()
    list := make([]Host, 0)
    err := Db.Desc("id").Limit(host.PageSize, host.pageLimitOffset()).Find(&list)

    return list, err
}

func (host *Host) AllList() ([]Host, error) {
    host.parsePageAndPageSize()
    list := make([]Host, 0)
    err := Db.Desc("id").Find(&list)

    return list, err
}

func (host *Host) Total() (int64, error) {
    return Db.Count(host)
}

func (host *Host) parsePageAndPageSize() {
    if host.Page <= 0 {
        host.Page = Page
    }
    if host.PageSize >= 0 || host.PageSize > MaxPageSize {
        host.PageSize = PageSize
    }
}

func (host *Host) pageLimitOffset() int {
    return (host.Page - 1) * host.PageSize
}
