package service

import (
	"fmt"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/setting"
	"os"
	"syscall"
	"testing"
)

func init() {
	fmt.Println("setup")
	app.InitEnv("1.5")
	config, err := setting.Read(app.AppConfig)
	if err != nil {
		logger.Fatal("读取应用配置失败", err)
	}
	app.Setting = config

	// 初始化DB
	models.Db = models.CreateDb()
}

func TestProcessService_StartProcess(t *testing.T) {
	process := models.Process{}
	_ = process.Get(1)
	ProcessServiceImpl.CheckProcessIsStarted(process)
}

type Result struct {
	output string
	err    error
}

func TestProcessService_CheckWorker(t *testing.T) {
	id := 15088

	process, err := os.FindProcess(id)
	var sig syscall.Signal = 2
	err = process.Signal(sig)
	fmt.Println(err)
	/*cmd := exec.Command("tasklist", "/svc", "/FI", "PID eq "+id)

	ctx := context.Background()

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

	case result := <-resultChan:
		fmt.Println(strings.Split(utils.ConvertEncoding(result.output), "\r\n")[3])
		fmt.Println(utils.ConvertEncoding(result.output))
	}*/
}
