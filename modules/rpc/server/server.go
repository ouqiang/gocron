package server

import (
    "golang.org/x/net/context"
    "net"
    "google.golang.org/grpc/grpclog"
    "google.golang.org/grpc"
    pb "github.com/ouqiang/gocron/modules/rpc/proto"
    "github.com/ouqiang/gocron/modules/utils"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/metadata"
    "errors"
)

type Server struct
{
    Token string
}

func (s Server) auth(ctx context.Context) error  {
    // 验证token是否有效
    meta, ok := metadata.FromContext(ctx)
    if !ok {
        return errors.New("missing metadata")
    }

    token, ok := meta["token"]
    if !ok {
        return errors.New("missing param token")
    }
    if token[0] != s.Token {
        return errors.New("invalid token")
    }

    return nil
}

func (s Server) Run(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error)  {
    defer func() {
        if err := recover(); err != nil {
            grpclog.Println(err)
        }
    } ()


    if s.Token != "" {
        err := s.auth(ctx)
        if err != nil {
            return nil, err
        }
    }

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

func Start(addr, certFile, keyFile, token string)  {
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
    server := Server{Token: token}
    if certFile != "" {
        // TLS认证
        creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
        if err != nil {
            grpclog.Fatalf("Failed to generate credentials %v", err)
        }

        s = grpc.NewServer(grpc.Creds(creds))
        pb.RegisterTaskServer(s, server)
        grpclog.Printf("listen %s with TLS", addr)
    } else {
        s = grpc.NewServer()
        pb.RegisterTaskServer(s, server)
        grpclog.Println("listen ", addr)
    }
    err = s.Serve(l)
    if err != nil {
        grpclog.Fatal(err)
    }
}

