package models

import "errors"

// 创建数据库表

type Migration struct{}

func (migration *Migration) Exec(dbName string) error {
	if !isDatabaseExist(dbName) {
		return errors.New("数据库不存在")
	}
	tables := []interface{}{
		&User{}, &Task{}, &TaskLog{}, &Host{},
	}
	for _, table := range tables {
		err := Db.Sync2(table)
		if err != nil {
			return err
		}
	}

	return nil
}

// 创建数据库
func isDatabaseExist(name string) bool {
	_, err := Db.Exec("use ?", name)

	return err != nil
}
