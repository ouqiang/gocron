package server

import (
	"net"
	"time"

	"google.golang.org/grpc/keepalive"

	"github.com/ouqiang/gocron/internal/modules/rpc/auth"
	pb "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

type Server struct{}

var keepAlivePolicy = keepalive.EnforcementPolicy{
	MinTime:             10 * time.Second,
	PermitWithoutStream: true,
}

var keepAliveParams = keepalive.ServerParameters{
	MaxConnectionIdle:     1 * time.Minute,
	MaxConnectionAge:      2 * time.Hour,
	MaxConnectionAgeGrace: 3 * time.Hour,
	Time:    30 * time.Second,
	Timeout: 3 * time.Second,
}

func (s Server) Run(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			grpclog.Println(err)
		}
	}()
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

func Start(addr string, enableTLS bool, certificate auth.Certificate) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		grpclog.Fatal(err)
	}

	var s *grpc.Server
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepAliveParams),
		grpc.KeepaliveEnforcementPolicy(keepAlivePolicy),
	}
	if enableTLS {
		tlsConfig, err := certificate.GetTLSConfigForServer()
		if err != nil {
			grpclog.Fatal(err)
		}
		opt := grpc.Creds(credentials.NewTLS(tlsConfig))
		opts = append(opts, opt)
	}
	s = grpc.NewServer(opts...)
	pb.RegisterTaskServer(s, Server{})
	grpclog.Printf("server listen on %s", addr)

	err = s.Serve(l)
	if err != nil {
		grpclog.Fatal(err)
	}
}
