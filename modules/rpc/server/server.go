package server

import (
    "golang.org/x/net/context"
    "net"
    "google.golang.org/grpc/grpclog"
    "google.golang.org/grpc"
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "github.com/ouqiang/gocron/modules/utils"
)

type Server struct {}

func (s Server) Run(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error)  {
    output, err := utils.ExecShell(ctx, req.Command)
    resp := new(pb.TaskResponse)
    resp.Output = output
    if err != nil {
        resp.Error = err.Error()
    } else {
        resp.Error = ""
    }

    return resp, nil
}

func Start(addr string)  {
    l, err := net.Listen("tcp", addr)
    if err != nil {
        grpclog.Fatal(err)
    }
    s := grpc.NewServer()
    pb.RegisterTaskServer(s, Server{})
    grpclog.Println("listen address ", addr)
    err = s.Serve(l)
    if err != nil {
        grpclog.Fatal(err)
    }
}

