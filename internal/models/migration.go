package models

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/ouqiang/gocron/internal/modules/logger"
)

type Migration struct{}

// 首次安装, 创建数据库表
func (migration *Migration) Install(dbName string) error {
	setting := new(Setting)
	task := new(Task)
	tables := []interface{}{
		&User{}, task, &TaskLog{}, &Host{}, setting, &LoginLog{}, &TaskHost{},
	}
	for _, table := range tables {
		exist, err := Db.IsTableExist(table)
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

	return nil
}

// 迭代升级数据库, 新建表、新增字段等
func (migration *Migration) Upgrade(oldVersionId int) {
	// v1.2版本不支持升级
	if oldVersionId == 120 {
		return
	}

	versionIds := []int{110, 122, 130, 140, 150}
	upgradeFuncs := []func(*xorm.Session) error{
		migration.upgradeFor110,
		migration.upgradeFor122,
		migration.upgradeFor130,
		migration.upgradeFor140,
		migration.upgradeFor150,
	}

	startIndex := -1
	// 从当前版本的下一版本开始升级
	for i, value := range versionIds {
		if value > oldVersionId {
			startIndex = i
			break
		}
	}

	if startIndex == -1 {
		return
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
			logger.Fatalf("事务回滚失败-%s", dbErr.Error())
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

// 升级到1.2.2版本
func (migration *Migration) upgradeFor122(session *xorm.Session) error {
	logger.Info("开始升级到v1.2.2")

	tableName := TablePrefix + "task"
	// task表增加tag字段
	_, err := session.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN tag VARCHAR(32) NOT NULL DEFAULT '' ", tableName))

	logger.Info("已升级到v1.2.2\n")

	return err
}

// 升级到v1.3版本
func (migration *Migration) upgradeFor130(session *xorm.Session) error {
	logger.Info("开始升级到v1.3")

	tableName := TablePrefix + "user"
	// 删除user表deleted字段
	_, err := session.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN deleted", tableName))

	logger.Info("已升级到v1.3\n")

	return err
}

// 升级到v1.4版本
func (migration *Migration) upgradeFor140(session *xorm.Session) error {
	logger.Info("开始升级到v1.4")

	tableName := TablePrefix + "task"
	// task表增加字段
	// retry_interval 重试间隔时间(秒)
	// http_method    http请求方法
	sql := fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN retry_interval SMALLINT NOT NULL DEFAULT 0,ADD COLUMN http_method TINYINT NOT NULL DEFAULT 1", tableName)
	_, err := session.Exec(sql)

	if err != nil {
		return err
	}

	logger.Info("已升级到v1.4\n")

	return err
}

func (m *Migration) upgradeFor150(session *xorm.Session) error {
	logger.Info("开始升级到v1.5")

	tableName := TablePrefix + "task"
	// task表增加字段 notify_keyword
	sql := fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN notify_keyword VARCHAR(128) NOT NULL DEFAULT '' ", tableName)
	_, err := session.Exec(sql)

	if err != nil {
		return err
	}

	settingModel := new(Setting)
	settingModel.Code = MailCode
	settingModel.Key = MailTemplateKey
	settingModel.Value = emailTemplate
	_, err = Db.Insert(settingModel)
	if err != nil {
		return err
	}
	settingModel.Id = 0
	settingModel.Code = SlackCode
	settingModel.Key = SlackTemplateKey
	settingModel.Value = slackTemplate
	_, err = Db.Insert(settingModel)
	if err != nil {
		return err
	}

	settingModel.Id = 0
	settingModel.Code = WebhookCode
	settingModel.Key = WebhookUrlKey
	settingModel.Value = ""
	_, err = Db.Insert(settingModel)
	if err != nil {
		return err
	}

	settingModel.Id = 0
	settingModel.Code = WebhookCode
	settingModel.Key = WebhookTemplateKey
	settingModel.Value = webhookTemplate
	_, err = Db.Insert(settingModel)
	if err != nil {
		return err
	}

	logger.Info("已升级到v1.5\n")

	return nil
}
