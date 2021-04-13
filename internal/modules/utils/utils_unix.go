// +build !windows

package utils

import (
	"errors"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"os/exec"
	"sync/atomic"
	"syscall"

	"golang.org/x/net/context"
)

type Result struct {
	output string
	err    error
}

// 执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	var flag atomic.Value
	resultChan := make(chan Result)
	go func() {
		output, err := cmd.CombinedOutput()
		if flag.Load() != nil && flag.Load().(bool) {
			return
		}
		resultChan <- Result{string(output), err}
	}()
	select {
	case <-ctx.Done():
		flag.Store(true)
		if err := cmd.Process.Kill(); err != nil {
			logger.Errorf("Process kill pid:%d err:%v", cmd.Process.Pid, err)
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return result.output, result.err
	}
}
