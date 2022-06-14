//go:build !windows
// +build !windows

package utils

import (
	"errors"
	"fmt"
	rpc "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"golang.org/x/net/context"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
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

// 启动一个worker进程
func StartWorker(ctx context.Context, req *rpc.StartRequest) (int, error) {
	cmd := exec.Command("/bin/bash", "-c", req.Command)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if req.LogFile != "" {
		logFile := req.LogFile
		d, err := filepath.Abs(path.Dir(logFile))
		if err != nil {
			return 0, err
		}
		_, err = os.Stat(d)
		if err != nil || os.IsNotExist(err) {
			_ = os.MkdirAll(d, 0666)
		}
		stdout, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(os.Getpid(), ": 打开日志文件错误:", err)
			return 0, err
		}
		cmd.Stderr = stdout
		cmd.Stdout = stdout
	}

	err := cmd.Start()

	pid := cmd.Process.Pid
	return pid, err
}

// 通过 pid 停止指定进程
func StopWorker(pid int64) error {
	//_ = exec.Command("kill", "-9", strconv.Itoa(int(pid))).Run()
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}
	return process.Kill()
}

func WorkerStateCheck(pid int64) (string, error) {
	if pid == 0 {
		return Error, errors.New("pid 不能为0")
	}
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Stop, err
	}
	rows := strings.Split(string(output), "\n")
	if len(rows) < 2 { //异常
		return Stop, errors.New(fmt.Sprintf("output exception,process %d not found", pid))
	}
	if strings.Contains(rows[1], "<defunct>") {
		//僵尸进程
		return Stop, errors.New(fmt.Sprintf("process %d is zombie", pid))
	}
	if !strings.Contains(rows[1], fmt.Sprintf("%d", pid)) {
		return Stop, errors.New(fmt.Sprintf("process %d not found", pid))
	}
	return Running, nil
}
