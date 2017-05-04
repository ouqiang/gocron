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
func Write(config map[string]string, filename string) error {
    if len(config) == 0 {
        return errors.New("参数不能为空")
    }

    file := ini.Empty()

    section, err := file.NewSection(DefaultSection)
    if err != nil {
        return err
    }
    for key, value := range config {
        if key == "" {
            continue
        }
        _, err = section.NewKey(key, value)
        if err != nil {
            return err
        }
    }
    err = file.SaveTo(filename)

    return err
}
