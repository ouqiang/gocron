// +build windows

package utils

import (
    "syscall"
    "time"
    "os/exec"
    "strconv"
    "fmt"
)

// 执行shell命令，可设置执行超时时间
func ExecShellWithTimeout(timeout int, command string, args... string) (string, error)  {
    cmd := exec.Command(command, args...)
    // 隐藏cmd窗口
    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
    // 后台运行
    if timeout == -1 {
        go cmd.CombinedOutput()
        return "", nil
    }
    // 不限制超时
    if timeout <= 0 {
        output ,err := cmd.CombinedOutput()
        return string(output), err
    }

    d := time.Duration(timeout) * time.Second
    timer := time.AfterFunc(d, func() {
        // 超时kill进程
        exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
        cmd.Process.Kill()
    })
    output ,err := cmd.CombinedOutput()
    timer.Stop()

    return string(output), err
}

// 格式化环境变量
func FormatEnv(key, value string) string {
    return fmt.Sprintf("set %s=%s & ", key, value)
}