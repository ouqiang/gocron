package models

import (
    "errors"
)


type Migration struct{}

// 首次安装, 创建数据库表
func (migration *Migration) Install(dbName string) error {
    if !isDatabaseExist(dbName) {
        return errors.New("数据库不存在")
    }
    setting := new(Setting)
    task := new(Task)
    tables := []interface{}{
        &User{}, task, &TaskLog{}, &Host{}, setting,&LoginLog{},&TaskHost{},
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
    setting.InitBasicField()
    task.CreateTestTask()

    return nil
}

// 判断数据库是否存在
func isDatabaseExist(name string) bool {
    _, err := Db.Exec("use ?", name)

    return err != nil
}

// 迭代升级数据库, 新建表、新增字段等
func (migration *Migration) Upgrade(oldVersionId int)  {
    versionIds   := []int{}
    upgradeFuncs := []func(){}

    startIndex := 0
    for i, value := range versionIds {
        if oldVersionId == value {
            startIndex = i + 1
            break;
        }
    }

    length := len(versionIds)
    for startIndex < length {
        upgradeFuncs[startIndex]()
        startIndex++
    }
}
