// +build windows

package utils

import (
    "syscall"
    "time"
    "os/exec"
    "strconv"
)

// 执行shell命令，可设置执行超时时间
func ExecShellWithTimeout(timeout int, command string) (string, error)  {
    cmd := exec.Command("cmd", "/C", command)
    // 隐藏cmd窗口
    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
    // 不限制超时
    if timeout <= 0 {
        output ,err := cmd.CombinedOutput()

        return ConvertEncoding(string(output), err)
    }

    d := time.Duration(timeout) * time.Second
    timer := time.AfterFunc(d, func() {
        // 超时kill进程
        exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
        cmd.Process.Kill()
    })
    output ,err := cmd.CombinedOutput()
    timer.Stop()

    return ConvertEncoding(string(output), err)
}

func ConvertEncoding(outputGBK string, err error) (string, error) {
    // windows平台编码为gbk，需转换为utf8才能入库
    outputUTF8, ok := GBK2UTF8(outputGBK)
    if ok {
        return outputUTF8, err
    }

    return "命令输出转换编码失败(gbk to utf8)", err
}