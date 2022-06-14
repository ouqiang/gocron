package grpcpool

import (
	"fmt"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/logger"
	rpc "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"github.com/ouqiang/gocron/internal/modules/setting"
	"golang.org/x/net/context"
	"testing"
)

func init() {
	fmt.Println("setup")
	app.InitEnv("1.5")
	config, err := setting.Read(app.AppConfig)
	if err != nil {
		logger.Fatal("读取应用配置失败", err)
	}
	app.Setting = config
}

func TestGRPCPool_Get(t *testing.T) {
	p := new(GRPCPool)
	p.conns = map[string]*Client{}
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 5921)
	client, err := p.Get(addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(client.Run(context.Background(), &rpc.TaskRequest{Command: "php -v"}))
	/*conn, err := grpc.DialContext(context.Background(), ":5921", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	client := rpc.NewTaskClient(conn)

	resp, err := client.Run(context.Background(), &rpc.TaskRequest{Command: "php -v"})
	t.Log(resp)*/
}
