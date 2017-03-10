package app

import (
	"os"
	"scheduler/utils"
	"runtime"
)

var  (
	AppDir string    // 应用根目录
	ConfDir string   // 配置目录
	LogDir string    // 日志目录
	DataDir string   // 数据目录，存放session文件等
	AppConfig string // 应用配置文件
	Installed bool = isInstalled()  // 应用是否安装过
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
}


// 检测环境
func CheckEnv()  {
	// ansible不支持安装在windows上, windows只能作为被控机
	if runtime.GOOS == "windows" {
		panic("不支持在windows上运行")
	}
	_, err := utils.ExecShell("ansible", "--version")
	if err != nil {
		panic(err)
	}
	_, err = utils.ExecShell("ansible-playbook", "--version")
	if err != nil {
		panic("ansible-playbook not found")
	}
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