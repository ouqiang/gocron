// +build !windows

package utils

import (
	"strings"
	"errors"
	"os/exec"
	"path/filepath"
	"syscall"
	"golang.org/x/net/context"
	"golang.org/x/exp/slices"
)

var langItems = []string{"python3", "python", "pwsh"}
type Result struct {
	output string
	err    error
}

// 执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error) {
	scriptArr := strings.Split(command, " ")
	cmd := exec.Command("/bin/bash", "-c", command)
	if(slices.Contains(langItems, scriptArr[0]) && len(scriptArr) > 1){
		cmd.Dir = filepath.Dir(scriptArr[1])
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	resultChan := make(chan Result)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- Result{string(output), err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return result.output, result.err
	}
}

