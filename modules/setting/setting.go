package setting

import (
    "errors"
    "gopkg.in/ini.v1"
)

const DefaultSection = "default"

// 读取配置
func Read(filename string) (*ini.Section,error) {
    config, err := ini.Load(filename)
    if err != nil {
        return nil, err
    }
    section := config.Section(DefaultSection)

    return section, nil
}

// 写入配置
func Write(config []string, filename string) error {
    if len(config) == 0 {
        return errors.New("参数不能为空")
    }
    if len(config) % 2 != 0 {
        return errors.New("参数不匹配")
    }

    file := ini.Empty()

    section, err := file.NewSection(DefaultSection)
    if err != nil {
        return err
    }
    for i := 0 ;i < len(config); {
        _, err = section.NewKey(config[i], config[i+1])
        if err != nil {
            return err
        }
        i += 2
    }
    err = file.SaveTo(filename)

    return err
}
