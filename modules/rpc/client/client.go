package client

import (
    "google.golang.org/grpc"
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "golang.org/x/net/context"
    "fmt"
)

func Exec(ip string, port int, taskReq *pb.TaskRequest) (output string, err error)  {
    addr := fmt.Sprintf("%s:%d", ip, port);
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return
    }
    defer conn.Close()
    c := pb.NewTaskClient(conn)
    resp, err := c.Run(context.Background(), taskReq)
    if err != nil {
        return
    }
    output = resp.Output

    return
}