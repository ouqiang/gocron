package models

// 创建数据库表

type Migration struct {}

func(migration *Migration) Exec() error  {
	tables := []interface{}{
		&User{}, &Task{}, &TaskLog{},&Host{},
	}
	for _, table := range(tables) {
		err := Db.Sync2(table)
		if err != nil {
			return err
		}
	}

	return nil
}