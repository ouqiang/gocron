package models

import (
    "github.com/ouqiang/gocron/modules/ssh"
    "github.com/go-xorm/xorm"
    "github.com/ouqiang/gocron/modules/app"
    "github.com/ouqiang/gocron/modules/utils"
    "errors"
    "io/ioutil"
    "strings"
    "github.com/ouqiang/gocron/modules/logger"
)

// 主机
type Host struct {
    Id        int16     `xorm:"smallint pk autoincr"`
    Name      string    `xorm:"varchar(64) notnull"`             // 主机名称
    Alias     string    `xorm:"varchar(32) notnull default '' "`  // 主机别名
    Username  string    `xorm:"varchar(32) notnull default '' "`  // ssh 用户名
    Port      int       `xorm:"notnull default 22"`               // 主机端口
    Remark    string    `xorm:"varchar(100) notnull default '' "` // 备注
    AuthType  ssh.HostAuthType      `xorm:"tinyint notnull default 1"` // 认证方式 1: 密码 2: 公钥
    BaseModel       `xorm:"-"`
}

// 新增
func (host *Host) Create() (insertId int16, err error) {
    _, err = Db.Insert(host)
    if err == nil {
        insertId = host.Id
    }

    return
}

func (host *Host) UpdateBean(id int16) (int64, error)  {
    return Db.ID(id).Cols("name,alias,username,port,remark,auth_type").Update(host)
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

func (host *Host) List(params CommonMap) ([]Host, error) {
    host.parsePageAndPageSize(params)
    list := make([]Host, 0)
    session := Db.Desc("id")
    host.parseWhere(session, params)
    err := session.Limit(host.PageSize, host.pageLimitOffset()).Find(&list)

    return list, err
}

func (host *Host) AllList() ([]Host, error) {
    list := make([]Host, 0)
    err := Db.Desc("id").Find(&list)

    return list, err
}

func (host *Host) Total(params CommonMap) (int64, error) {
    session := Db.NewSession()
    host.parseWhere(session, params)
    return session.Count(host)
}

func (h *Host) GetPasswordByHost(host string) (string, error) {
    path := app.DataDir + "/ssh/password/" + host


    return h.readFile(path)
}

func (h *Host) GetPrivateKeyByHost(host string) (string,error)  {
    path := app.DataDir + "/ssh/private_key/" + host

    return h.readFile(path)
}

func (host *Host) readFile(file string) (string, error) {
    logger.Debug("认证文件路径: ", file)
    if !utils.FileExist(file) {
        return "", errors.New(file + "-认证文件不存在或无权限访问")
    }

    contentByte, err := ioutil.ReadFile(file)
    if err != nil {
        return "", err
    }
    content := string(contentByte)
    content = strings.TrimSpace(content)
    if content == "" {
        return "", errors.New("密码为空")
    }

    return content, nil
}

// 解析where
func (host *Host) parseWhere(session *xorm.Session, params CommonMap)  {
    if len(params) == 0 {
        return
    }
    id, ok := params["Id"]
    if ok && id.(int) > 0 {
        session.And("id = ?", id)
    }
    name, ok := params["Name"]
    if ok && name.(string) != "" {
        session.And("name = ?", name)
    }
}