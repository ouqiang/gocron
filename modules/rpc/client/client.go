package client

import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/grpclog"
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "golang.org/x/net/context"
)

func Start()  {
    conn, err := grpc.Dial("127.0.0.1:50000", grpc.WithInsecure())
    if err != nil {
        grpclog.Fatal(err)
    }
    defer conn.Close()
    c := pb.NewTaskClient(conn)
    req := new(pb.TaskRequest)
    resp, err := c.Run(context.Background(), req)
    if err != nil {
        grpclog.Fatal(err)
    }
    grpclog.Println(resp.Name)
}