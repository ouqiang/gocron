package client

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/rpc/grpcpool"
	pb "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	taskMap sync.Map
)

var (
	errUnavailable = errors.New("无法连接远程服务器")
)

func generateTaskUniqueKey(ip string, port int, id int64) string {
	return fmt.Sprintf("%s:%d:%d", ip, port, id)
}

func Stop(ip string, port int, id int64) {
	key := generateTaskUniqueKey(ip, port, id)
	cancel, ok := taskMap.Load(key)
	if !ok {
		return
	}
	cancel.(context.CancelFunc)()
}

func Exec(ip string, port int, taskReq *pb.TaskRequest) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("panic#rpc/client.go:Exec#", err)
		}
	}()
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := grpcpool.Pool.Get(addr)
	if err != nil {
		return "", err
	}
	isConnClosed := false
	defer func() {
		if !isConnClosed {
			grpcpool.Pool.Put(addr, conn)
		}
	}()
	c := pb.NewTaskClient(conn)
	if taskReq.Timeout <= 0 || taskReq.Timeout > 86400 {
		taskReq.Timeout = 86400
	}
	timeout := time.Duration(taskReq.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	taskUniqueKey := generateTaskUniqueKey(ip, port, taskReq.Id)
	taskMap.Store(taskUniqueKey, cancel)
	defer taskMap.Delete(taskUniqueKey)
	resp, err := c.Run(ctx, taskReq)
	if err != nil {
		return parseGRPCError(err, conn, &isConnClosed)
	}

	if resp.Error == "" {
		return resp.Output, nil
	}

	return resp.Output, errors.New(resp.Error)
}

func parseGRPCError(err error, conn *grpc.ClientConn, connClosed *bool) (string, error) {
	switch grpc.Code(err) {
	case codes.Unavailable, codes.Internal:
		conn.Close()
		*connClosed = true
		return "", errUnavailable
	case codes.DeadlineExceeded:
		return "", errors.New("执行超时, 强制结束")
	case codes.Canceled:
		return "", errors.New("手动停止")
	}
	return "", err
}
