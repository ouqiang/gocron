package server

import (
    "golang.org/x/net/context"
    "net"
    "google.golang.org/grpc/grpclog"
    "google.golang.org/grpc"
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
)

type Server struct {}

func (s Server) Run(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error)  {
    resp := new(pb.TaskResponse)
    resp.Name = "gRPC"

    return resp, nil
}

func Start()  {
    l, err := net.Listen("tcp", "127.0.0.1:50000")
    if err != nil {
        grpclog.Fatal(err)
    }
    s := grpc.NewServer()
    pb.RegisterTaskServer(s, Server{})
    grpclog.Println("listen address ", "127.0.0.1:50000")
    s.Serve(l)
}

