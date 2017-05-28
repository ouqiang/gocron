package client

import (
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "golang.org/x/net/context"
    "fmt"
    "time"
    "errors"
    "github.com/ouqiang/gocron/modules/rpc/grpcpool"
)

func Exec(ip string, port int, taskReq *pb.TaskRequest) (string, error)  {
    addr := fmt.Sprintf("%s:%d", ip, port);
    conn, err := grpcpool.Pool.Get(addr)
    if err != nil {
        return "", err
    }
    defer func() {
        grpcpool.Pool.Put(addr, conn)
    }()
    c := pb.NewTaskClient(conn)
    if taskReq.Timeout <= 0 || taskReq.Timeout > 86400 {
        taskReq.Timeout = 86400
    }
    timeout := time.Duration(taskReq.Timeout) * time.Second
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    resp, err := c.Run(ctx, taskReq)
    if err != nil {
        return "", err
    }

    if resp.Error == "" {
        return resp.Output, nil
    }

    return resp.Output, errors.New(resp.Error)
}