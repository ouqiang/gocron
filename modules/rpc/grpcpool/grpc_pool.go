package grpcpool

import (
    "github.com/silenceper/pool"
    "sync"
    "time"
    "google.golang.org/grpc"
    "errors"
    "google.golang.org/grpc/credentials"
    "golang.org/x/net/context"
    "strings"
)


var (
    Pool GRPCPool
)

var (
    ErrInvalidConn = errors.New("invalid connection")
)

func init()  {
    Pool = GRPCPool{
        make(map[string]pool.Pool),
        sync.RWMutex{},
    }
}

type GRPCPool struct {
    // map key格式 ip:port
    conns map[string]pool.Pool
    sync.RWMutex
}

func (p *GRPCPool) Get(addr, certFile, token string) (*grpc.ClientConn, error)  {
    p.RLock()
    pool, ok := p.conns[addr]
    p.RUnlock()
    if !ok {
        err := p.newCommonPool(addr, certFile, token)
        if err != nil {
            return nil, err
        }
    }

    p.RLock()
    pool = p.conns[addr]
    p.RUnlock()
    conn, err := pool.Get()
    if err != nil {
        return nil, err
    }

    return conn.(*grpc.ClientConn), nil
}

func (p *GRPCPool) Put(addr string, conn *grpc.ClientConn) error {
    p.RLock()
    defer p.RUnlock()
    pool, ok := p.conns[addr]
    if ok {
        return pool.Put(conn)
    }

    return ErrInvalidConn
}


// 释放连接池
func (p *GRPCPool) Release(addr string) {
    p.Lock()
    defer p.Unlock()
    pool, ok := p.conns[addr]
    if !ok {
        return
    }
    pool.Release()
    delete(p.conns, addr)
}

// 释放所有连接池
func (p *GRPCPool) ReleaseAll()  {
    p.Lock()
    defer p.Unlock()
    for _, pool := range(p.conns) {
        pool.Release()
    }
}

// 初始化底层连接池
func (p *GRPCPool) newCommonPool(addr, certFile, token string) (error) {
    p.Lock()
    defer p.Unlock()
    commonPool, ok := p.conns[addr]
    if ok {
        return nil
    }
    poolConfig := &pool.PoolConfig{
        InitialCap: 1,
        MaxCap: 30,
        Factory: func() (interface{}, error) {
            if certFile == "" {
                return grpc.Dial(addr, grpc.WithInsecure())
            }

            server := strings.Split(addr, ":")
            creds, err := credentials.NewClientTLSFromFile(certFile, server[0])
            if err != nil {
                return nil, err
            }

            customCredential := &CustomCredential{Token: token}


            return grpc.Dial(addr,
                grpc.WithTransportCredentials(creds),
                grpc.WithPerRPCCredentials(customCredential),
            )
        },
        Close: func(v interface{}) error {
            conn, ok := v.(*grpc.ClientConn)
            if ok && conn != nil {
                return conn.Close()
            }
            return ErrInvalidConn
        },
        IdleTimeout: 3 * time.Minute,
    }

    commonPool, err := pool.NewChannelPool(poolConfig)
    if err != nil {
        return err
    }

    p.conns[addr] = commonPool

    return nil
}

type CustomCredential struct
{
    Token string
}

func (c CustomCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
    return map[string]string{
        "token": c.Token,
    }, nil
}

func (c CustomCredential) RequireTransportSecurity() bool {
    return true
}