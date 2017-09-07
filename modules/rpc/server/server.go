package server

import (
    "golang.org/x/net/context"
    "net"
    "google.golang.org/grpc/grpclog"
    "google.golang.org/grpc"
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "github.com/ouqiang/gocron/modules/utils"
    "google.golang.org/grpc/credentials"
)

type Server struct {}

func (s Server) Run(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error)  {
    defer func() {
        if err := recover(); err != nil {
            grpclog.Println(err)
        }
    } ()
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

func Start(addr, certFile, keyFile string)  {
    defer func() {
       if err := recover(); err != nil {
           grpclog.Println("panic", err)
       }
    } ()
    l, err := net.Listen("tcp", addr)
    if err != nil {
        grpclog.Fatal(err)
    }

    var s *grpc.Server
    if certFile != "" {
        // TLS认证
        creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
        if err != nil {
            grpclog.Fatalf("Failed to generate credentials %v", err)
        }

        s = grpc.NewServer(grpc.Creds(creds))
        pb.RegisterTaskServer(s, Server{})
        grpclog.Printf("listen %s with TLS", addr)
    } else {
        s = grpc.NewServer()
        pb.RegisterTaskServer(s, Server{})
        grpclog.Println("listen ", addr)
    }
    err = s.Serve(l)
    if err != nil {
        grpclog.Fatal(err)
    }
}

