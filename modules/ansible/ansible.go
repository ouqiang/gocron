package ansible

// ansible ad-hoc 命令封装

import (
    "errors"
    "github.com/ouqiang/gocron/modules/utils"
    "strings"
)

// ansible是否有安装
func IsInstalled() bool  {
    _, err := utils.ExecShell("ansible", "--version")


    return err == nil
}

/**
 * 执行ad-hoc
 * hosts  主机名或主机别名 逗号分隔
 * hostFile 主机名文件
 * module 模块
 * args   传递给module的参数
 */
func ExecCommand(hosts string, hostFile string, module string, args ...string) (output string, err error) {
    if hosts == "" || hostFile == "" || module == "" {
        err = errors.New("参数不完整")
        return
    }
    commandArgs := []string{hosts, "-i", hostFile, "-m", module}
    if len(args) > 0 {
        commandArgs = append(commandArgs, args...)
    }
    retryTimes := 2
    for i := 0; i < retryTimes; i++ {
        output, err = utils.ExecShell("ansible", commandArgs...)
        // todo 偶尔出现Failed to connect to the host via ssh. 暂时重试解决
        if err == nil || !strings.HasPrefix(err.Error(), "exit status 3") {
            break;
        }
    }


    return
}

// 执行shell命令
func Shell(hosts string, hostFile string, args ...string) (output string, err error) {
    return ExecCommand(hosts, hostFile, "shell", args...)
}

// 复制本地脚本到远程执行
func Script(hosts string, hostFile string, args ...string) (output string, err error) {
    return ExecCommand(hosts, hostFile, "script", args...)
}

// 测试主机是否可通
func Ping(hosts string, hostFile string) (output string, err error) {
    return ExecCommand(hosts, hostFile, "ping")
}
