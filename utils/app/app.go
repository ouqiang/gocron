package app

import (
	"os"
)

var  (
	AppDir string    // 应用根目录
	ConfDir string   // 配置目录
	LogDir string    // 日志目录
	DataDir string   // 数据目录，存放session文件等
	AppConfig string // 应用配置文件
	Installed bool   // 应用是否安装过
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	AppDir = wd
	ConfDir = AppDir + "/conf"
	LogDir  = AppDir + "/log"
	DataDir = AppDir + "/data"
	AppConfig = AppDir + "/app.ini"
	checkDirExists(ConfDir, LogDir, DataDir)
	Installed = isInstalled()
}

// 判断应用是否安装过
func isInstalled() bool {
	_, err := os.Stat(ConfDir + "/install.lock")
	if os.IsExist(err) {
		return true
	}

	return false
}

// 检测目录是否存在
func checkDirExists(path... string)  {
	for _, value := range(path) {
		_, err := os.Stat(value)
		if os.IsNotExist(err) {
			panic(value + "目录不存在")
		}
		if os.IsPermission(err) {
			panic(value + "目录无权限操作")
		}
	}
}