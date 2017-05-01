package app

import (
    "os"

    "github.com/ouqiang/gocron/modules/logger"
    "runtime"
    "github.com/ouqiang/gocron/modules/utils"
)

var (
    AppDir       string // 应用根目录
    ConfDir      string // 配置目录
    LogDir       string // 日志目录
    DataDir      string // 存放session等
    AppConfig    string // 应用配置文件
    Installed    bool   // 应用是否安装过
)

func InitEnv() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    logger.InitLogger()
    wd, err := os.Getwd()
    if err != nil {
        logger.Fatal(err)
    }
    AppDir = wd
    ConfDir = AppDir + "/conf"
    LogDir = AppDir + "/log"
    DataDir = AppDir + "/data"
    AppConfig = ConfDir + "/app.ini"
    checkDirExists(ConfDir, LogDir, DataDir)
    Installed = IsInstalled()
}

// 判断应用是否安装过
func IsInstalled() bool {
    _, err := os.Stat(ConfDir + "/install.lock")
    if os.IsNotExist(err) {
        return false
    }

    return true
}

// 创建安装锁文件
func CreateInstallLock() error {
    _, err := os.Create(ConfDir + "/install.lock")
    if err != nil {
        logger.Error("创建安装锁文件失败")
    }

    return err
}

// 检测目录是否存在
func checkDirExists(path ...string) {
    for _, value := range path {
        if !utils.FileExist(value) {
            logger.Fatal(value + "目录不存在或无权限访问")
        }
    }
}