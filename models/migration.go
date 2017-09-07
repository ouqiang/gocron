package models

import (
    "errors"
    "fmt"
    "github.com/ouqiang/gocron/modules/logger"
    "github.com/go-xorm/xorm"
    "strconv"
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
    versionIds   := []int{110}
    upgradeFuncs := []func(*xorm.Session) error {
        migration.upgradeFor110,
    }

    // 默认当前版本为v1.0
    startIndex := 0
    // 从当前版本的下一版本开始升级
    for i, value := range versionIds {
        if oldVersionId == value {
            startIndex = i + 1
            break;
        }
    }

    length := len(versionIds)
    if startIndex >= length {
        return
    }

    session := Db.NewSession()
    err := session.Begin()
    if err != nil {
        logger.Fatalf("开启事务失败-%s", err.Error())
    }
    for startIndex < length {
        err = upgradeFuncs[startIndex](session)
        if err == nil {
            startIndex++
            continue
        }
        dbErr := session.Rollback()
        if dbErr != nil {
            logger.Fatalf("事务回滚失败-%s",dbErr.Error())
        }
        logger.Fatal(err)
    }
    err = session.Commit()
    if err != nil {
        logger.Fatalf("提交事务失败-%s", err.Error())
    }
}

// 升级到v1.1版本
func (migration *Migration) upgradeFor110(session *xorm.Session) error {
    logger.Info("开始升级到v1.1")
    // 创建表task_host
    err := session.Sync2(new(TaskHost))
    if err != nil {
        return err
    }

    tableName := TablePrefix + "task"
    // 把task对应的host_id写入task_host表
    sql := fmt.Sprintf("SELECT id, host_id FROM %s WHERE host_id > 0", tableName)
    results, err := session.Query(sql)
    if err != nil {
        return err
    }

    for _, value := range results {
        taskHostModel := &TaskHost{}
        taskId, err := strconv.Atoi(string(value["id"]))
        if err != nil {
            return err
        }
        hostId, err := strconv.Atoi(string(value["host_id"]))
        if err != nil {
            return err
        }
        taskHostModel.TaskId = taskId
        taskHostModel.HostId = int16(hostId)
        _, err = session.Insert(taskHostModel)
        if err != nil {
            return err
        }
    }


    // 删除task表host_id字段
    _, err = session.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN host_id", tableName))

    logger.Info("已升级到v1.1\n")

    return err
}
