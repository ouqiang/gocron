package client

import (
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "golang.org/x/net/context"
    "fmt"
    "time"
    "errors"
    "github.com/ouqiang/gocron/modules/rpc/grpcpool"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc"
)

func Exec(ip string, port int, taskReq *pb.TaskRequest) (string, error)  {
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
    defer cancel()
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
        case codes.Unavailable:
            conn.Close()
            *connClosed = true
            return "", errors.New("无法连接远程服务器")
        case codes.DeadlineExceeded:
            return "", errors.New("执行超时, 强制结束")
    }
    return "", err
}
