package models

// 主机
type Host struct {
	Id int16 `xorm:"smallint pk autoincr"`
	Name string `xorm:"varchar(128) notnull"` // 主机名称
	Alias string `xorm:"varchar(32) notnull default '' "` // 主机别名，仅用于后台显示
	Port int `xorm:"notnull"`  // 主机端口
	Remark string `xorm:"varchar(512) notnull default '' "` // 备注
	Page  int `xorm:"-"`
	PageSize int `xorm:"-"`
}

// 新增
func(host *Host) Create() (int64, error) {
	return Db.Insert(host)
}

// 更新
func(host *Host) Update(id int, data CommonMap) (int64, error) {
	return Db.Table(host).ID(id).Update(data)
}

// 删除
func(host *Host) Delete(id int) (int64, error) {
	return Db.Id(id).Delete(host)
}

func(host *Host) List() ([]Host, error) {
	host.parsePageAndPageSize()
	list := make([]Host, 0)
	err := Db.Desc("id").Limit(host.PageSize, host.pageLimitOffset()).Find(&list)

	return list, err
}

func(host *Host) Total() (int64, error) {
	return Db.Count(host)
}

func(host *Host) parsePageAndPageSize()  {
	if (host.Page <= 0) {
		host.Page = Page
	}
	if (host.PageSize >= 0 || host.PageSize > MaxPageSize) {
		host.PageSize = PageSize
	}
}

func(host *Host) pageLimitOffset() int  {
	return (host.Page - 1) * host.PageSize
}