// +build !windows

package utils

import (
    "os/exec"
    "syscall"
    "golang.org/x/net/context"
    "errors"
)

type Result struct {
    output string
    err error
}

// 执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error)  {
    cmd := exec.Command("/bin/bash", "-c", command)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Setpgid: true,
    }
    var resultChan chan Result = make(chan Result)
    go func() {
        output ,err := cmd.CombinedOutput()
        resultChan <- Result{string(output), err}
    }()
    select {
        case <- ctx.Done():
            syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
            return "", errors.New("timeout killed")
        case result := <- resultChan:
            return result.output, result.err
    }


    return "", nil
}