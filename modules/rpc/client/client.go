package client

import (
    "google.golang.org/grpc"
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "golang.org/x/net/context"
    "fmt"
    "time"
)

func Exec(ip string, port int, taskReq *pb.TaskRequest) (output string, err error)  {
    addr := fmt.Sprintf("%s:%d", ip, port);
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return
    }
    defer conn.Close()
    c := pb.NewTaskClient(conn)
    if taskReq.Timeout <= 0 || taskReq.Timeout > 86400 {
        taskReq.Timeout = 86400
    }
    timeout := time.Duration(taskReq.Timeout) * time.Second
    ctx, _ := context.WithTimeout(context.Background(), timeout)
    resp, err := c.Run(ctx, taskReq)
    if err != nil {
        return
    }
    output = resp.Output

    return
}