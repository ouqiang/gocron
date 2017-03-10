package setting

import (
	"gopkg.in/ini.v1"
	"errors"
	"scheduler/utils/app"
)

// 读取配置
func Read() (config *ini.File, err error) {
	config, err = ini.Load(app.AppConfig)
	if err != nil {
		return
	}

	return
}


// 写入配置
func Write(config map[string]map[string]string) (error) {
	if len(config) == 0 {
		return errors.New("参数不能为空")
	}

	file := ini.Empty()
	for sectionName, items := range(config) {
		if sectionName == "" {
			return errors.New("节名称不能为空")
		}
		section, err := file.NewSection(sectionName)
		if err != nil {
			return err
		}
		for key, value := range(items) {
			_, err = section.NewKey(key, value)
			if err != nil {
				return err
			}
		}
	}
	err := file.SaveTo(app.AppConfig)

	return err
}