package grpcpool

import (
	"context"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
	"sync"
	"time"

	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/rpc/auth"
	"github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	backOffBaseDelay = 3 * time.Second
	backOffMaxDelay  = 60 * time.Second
	dialTimeout      = 2 * time.Second
)

var (
	Pool = &GRPCPool{
		conns: make(map[string]*Client),
	}

	keepAliveParams = keepalive.ClientParameters{
		Time:                20 * time.Second,
		Timeout:             3 * time.Second,
		PermitWithoutStream: true,
	}
)

type Client struct {
	conn          *grpc.ClientConn
	rpcClient     rpc.TaskClient
	processClient rpc.ProcessClient
}

type GRPCPool struct {
	// map key格式 ip:port
	conns map[string]*Client
	mu    sync.RWMutex
}

func (p *GRPCPool) Get(addr string) (rpc.TaskClient, error) {
	p.mu.RLock()
	client, ok := p.conns[addr]
	p.mu.RUnlock()
	if ok {
		return client.rpcClient, nil
	}

	client, err := p.factory(addr)
	if err != nil {
		return nil, err
	}

	return client.rpcClient, nil
}

func (p *GRPCPool) GetClient(addr string) (rpc.ProcessClient, error) {
	p.mu.RLock()
	client, ok := p.conns[addr]
	p.mu.RUnlock()
	if ok {
		return client.processClient, nil
	}

	client, err := p.factory(addr)
	if err != nil {
		return nil, err
	}
	return client.processClient, nil
}

// Release 释放连接
func (p *GRPCPool) Release(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	client, ok := p.conns[addr]
	if !ok {
		return
	}
	delete(p.conns, addr)
	client.conn.Close()
}

// 创建连接
func (p *GRPCPool) factory(addr string) (*Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	client, ok := p.conns[addr]
	if ok {
		return client, nil
	}
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepAliveParams),
		//grpc.WithBackoffMaxDelay(backOffMaxDelay),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.Config{
			BaseDelay: backOffBaseDelay,
			MaxDelay:  backOffMaxDelay,
		}}),
	}

	if !app.Setting.EnableTLS {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		server := strings.Split(addr, ":")
		certificate := auth.Certificate{
			CAFile:     app.Setting.CAFile,
			CertFile:   app.Setting.CertFile,
			KeyFile:    app.Setting.KeyFile,
			ServerName: server[0],
		}

		transportCreds, err := certificate.GetTransportCredsForClient()
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(transportCreds))
	}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}

	client = &Client{
		conn:          conn,
		rpcClient:     rpc.NewTaskClient(conn),
		processClient: rpc.NewProcessClient(conn),
	}

	p.conns[addr] = client

	return client, nil
}
