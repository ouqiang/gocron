package grpcpool

import (
    "github.com/silenceper/pool"
    "sync"
    "time"
    "google.golang.org/grpc"
    "errors"
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

func (p *GRPCPool) Get(addr string) (*grpc.ClientConn, error)  {
    p.RLock()
    p.RUnlock()
    pool, ok := p.conns[addr]
    if !ok {
        err := p.newCommonPool(addr)
        if err != nil {
            return nil, err
        }
    }

    pool = p.conns[addr]
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
func (p *GRPCPool) newCommonPool(addr string) (error) {
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
            return grpc.Dial(addr, grpc.WithInsecure())
        },
        Close: func(v interface{}) error {
            conn, ok := v.(*grpc.ClientConn)
            if ok && conn != nil {
                return conn.Close()
            }
            return ErrInvalidConn
        },
        IdleTimeout: 5 * time.Minute,
    }

    commonPool, err := pool.NewChannelPool(poolConfig)
    if err != nil {
        return err
    }

    p.conns[addr] = commonPool

    return nil
}