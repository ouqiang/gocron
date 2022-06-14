//go:build windows

package utils

import (
	"errors"
	"fmt"
	rpc "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/net/context"
)

type Result struct {
	output string
	err    error
}

// ExecShell 执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error) {
	cmd := exec.Command("cmd", "/C", command)
	// 隐藏cmd窗口
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	var resultChan chan Result = make(chan Result)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- Result{string(output), err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {
			exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
			cmd.Process.Kill()
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return ConvertEncoding(result.output), result.err
	}
}

// StartWorker 开始工作进程
func StartWorker(ctx context.Context, req *rpc.StartRequest) (int, error) {
	cmd := exec.Command("cmd", "/C", req.Command)
	// 隐藏cmd窗口
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
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

// StopWorker 通过 pid 停止指定工作进程
func StopWorker(pid int64) error {
	state, _ := WorkerStateCheck(pid)
	if state != Running {
		return nil
	}
	//强制终止此进程及其启动的任何子进程
	cmd := exec.Command("taskkill", "/PID", strconv.FormatInt(pid, 10), "/T", "/F")
	return cmd.Run()
}

func WorkerStateCheck(pid int64) (string, error) {
	if pid == 0 {
		return Error, errors.New("pid 不能为0")
	}
	cmd := exec.Command("tasklist", "/svc", "/FI", fmt.Sprintf("PID eq %d", pid))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Error, err
	}
	str := ConvertEncoding(string(output))
	rows := strings.Split(str, "\n")
	if len(rows) < 4 { //异常
		return Stop, errors.New(str)
	}
	re, err := regexp.Compile(" (\\d+) ")
	if err != nil {
		return Stop, err
	}
	match := re.FindStringSubmatch(rows[3])
	if len(match) < 2 {
		return Stop, errors.New("程序停止运行")
	}
	return Running, nil
}

func ConvertEncoding(outputGBK string) string {
	// windows平台编码为gbk，需转换为utf8才能入库
	outputUTF8, ok := GBK2UTF8(outputGBK)
	if ok {
		return outputUTF8
	}

	return outputGBK
}
