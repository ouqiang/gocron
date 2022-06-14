package service

import (
	"fmt"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/rpc/grpcpool"
	rpc "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ProcessService struct{}

var ProcessServiceImpl ProcessService

// Initialize 从数据库中取出所有开启的进程,定时检测进程是否在启动中
func (p ProcessService) Initialize() {
	ticker := time.NewTicker(30 * time.Second)
	go func(t *time.Ticker) {
		for {
			<-t.C
			//检测所有开始状态的进程是否在运行中
			var processes []models.Process
			_ = models.Db.Where("status = ? AND enable = ?", models.ProcessStart, 1).Find(&processes)
			for _, process := range processes {
				go p.CheckProcessIsStarted(process)
			}

			//todo 检测所有停止状态的进程是否有在运行中的
			/*var workers []models.ProcessWorker
			_ = models.Db.Where("is_valid = ?", 1).Find(&workers)
			for _, worker := range workers {
				go checkWorkerState(worker)
			}*/
		}
	}(ticker)
}

func (p ProcessService) CheckProcessIsStarted(process models.Process) {
	var workers []models.ProcessWorker
	_ = models.Db.Where("process_id = ? AND is_valid = ?", process.Id, 1).Find(&workers)
	ph := models.ProcessHost{}
	hosts := ph.GetByProcess(process)
	if len(hosts) == 0 {
		return
	}
	index := 0
	for _, worker := range workers {
		host := hosts[index]
		if worker.HostId == 0 || worker.State == models.Stop {
			worker.HostId = host.Id
			_ = worker.Update()
			index++
			if index == len(hosts) {
				index = 0
			}
		}
		if worker.Pid == 0 {
			workerStart(worker)
			_ = worker.SetState(models.Start)
		} else {
			p.CheckWorkerIsRunning(worker)
		}
	}
}

func (p ProcessService) CheckWorkerIsRunning(worker models.ProcessWorker) {
	state, err := getWorkerState(worker)
	logger.Info(err, state, worker)
	if err != nil {
		return
	}
	if state != utils.Running {
		workerStart(worker)
		_ = worker.SetState(models.Start)
	}
}

func getWorkerState(worker models.ProcessWorker) (string, error) {
	host := models.Host{}
	err := host.Find(worker.HostId)
	if err != nil {
		return "", err
	}
	addr := fmt.Sprintf("%s:%d", host.Name, host.Port)
	client, _ := grpcpool.Pool.GetClient(addr)
	resp, err := client.WorkerStateCheck(context.Background(), &rpc.PidRequest{
		Pid: worker.Pid,
	})
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			grpcpool.Pool.Release(addr) // 链接不可用,释放链接
		}
		return "", err
	}
	return resp.Code, nil
}

func checkWorkerState(worker models.ProcessWorker) {
	host := models.Host{}
	err := host.Find(worker.HostId)
	if err != nil {
		logger.Error(err)
		return
	}
	addr := fmt.Sprintf("%s:%d", host.Name, host.Port)
	client, _ := grpcpool.Pool.GetClient(addr)
	resp, err := client.WorkerStateCheck(context.Background(), &rpc.PidRequest{
		Pid: worker.Pid,
	})
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			grpcpool.Pool.Release(addr) // 链接不可用,释放链接
		}
		return
	}

	switch worker.State {
	case models.Starting:
		if resp.Code == utils.Running {
			//已启动
			_ = worker.SetState(models.Start)
		}
		if resp.Code == utils.Stop {
			//启动失败
			_ = worker.SetState(models.Exited)
			workerStart(worker)
		}
	case models.Start:
		if resp.Code == utils.Stop {
			//已停止,需要重新启动
			workerStart(worker)
		}
	case models.Stop:
		if resp.Code == utils.Running {
			//异常
		}
	case models.Exited:
		if resp.Code != utils.Running {
			workerStart(worker)
		}
	}
}

func workerStart(worker models.ProcessWorker) {
	logger.Debug("workerStart running")
	host := models.Host{}
	err := host.Find(worker.HostId)
	if err != nil {
		logger.Debug("get worker fail", err)
		return
	}
	addr := fmt.Sprintf("%s:%d", host.Name, host.Port)
	client, err := grpcpool.Pool.GetClient(addr)

	process := models.Process{}
	process.Get(worker.ProcessId)

	req := rpc.StartRequest{
		Command: process.Command,
		LogFile: process.LogFile,
	}
	resp, err := client.StartWorker(context.Background(), &req)
	worker.State = models.Starting
	worker.Pid = resp.Pid
	_ = worker.Update()

	//5秒后确认该进程是否还在还在运行中
	time.AfterFunc(time.Second*5, func() {
		checkWorkerState(worker)
		logger.Debug(fmt.Sprintf("%d is running", worker.Pid))
	})
}

func (p ProcessService) StopWorker(worker models.ProcessWorker) {
	host := models.Host{}
	err := host.Find(worker.HostId)
	if err != nil {
		return
	}
	addr := fmt.Sprintf("%s:%d", host.Name, host.Port)
	client, err := grpcpool.Pool.GetClient(addr)
	req := rpc.PidRequest{
		Pid: worker.Pid,
	}
	resp, _ := client.StopWorker(context.Background(), &req)
	if resp.Code == "success" {
		_ = worker.SetState(models.Stop)
	}
}

func (p ProcessService) StopProcess(process models.Process) error {
	pw := models.ProcessWorker{}
	workers, _ := pw.GetByProcess(process)
	for _, worker := range workers {
		p.StopWorker(worker)
	}
	return nil
}
