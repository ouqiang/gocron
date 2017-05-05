// +build !windows

package utils

import (
    "os/exec"
    "syscall"
    "time"
    "fmt"
)

// 执行shell命令，可设置执行超时时间
func ExecShellWithTimeout(timeout int, command string, args... string) (string, error)  {
    cmd := exec.Command(command, args...)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Setpgid: true,
    }

    // 后台运行
    if timeout == -1 {
        go cmd.CombinedOutput()
        return "", nil
    }
    // 不限制超时
    if timeout == 0 {
        output ,err := cmd.CombinedOutput()
        return string(output), err
    }

    d := time.Duration(timeout) * time.Second
    timer := time.AfterFunc(d, func() {
        // 超时kill进程
        syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
    })
    output ,err := cmd.CombinedOutput()
    timer.Stop()

    return string(output), err
}

// 格式化环境变量
func FormatEnv(key, value string) string {
    return fmt.Sprintf("export %s=%s;", key, value)
}