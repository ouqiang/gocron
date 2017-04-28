package models

import (
    "errors"
)

// 创建数据库表

type Migration struct{}

func (migration *Migration) Exec(dbName string) error {
    if !isDatabaseExist(dbName) {
        return errors.New("数据库不存在")
    }
    setting := new(Setting)
    tables := []interface{}{
        &User{}, &Task{}, &TaskLog{}, &Host{}, setting,
    }
    for _, table := range tables {
        exist, err:= Db.IsTableExist(table)
        if exist {
            return errors.New("数据表已存在")
        }
        if err != nil {
            return err
        }
        err = Db.Sync2(table)
        if err != nil {
            return err
        }
    }
    err := setting.InitBasicField()
    if err != nil {
        return err
    }

    return nil
}

// 创建数据库
func isDatabaseExist(name string) bool {
    _, err := Db.Exec("use ?", name)

    return err != nil
}
