package main

import (
	"fmt"
	rpc "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	s := " php -r while(true){var_dump(date(\"Y-m-d H:i:s\"));};"
	cmd := exec.Command("cmd", "/C", "php -r 'while(true){var_dump(date(\"Y-m-d H:i:s\"));};'")
	fmt.Println(strings.Split(s, "/C"))
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(os.Getpid())
	t.Log(cmd.Process.Pid)

	time.Sleep(30 * time.Second)
}

func TestStartProcess(t *testing.T) {
	//connection rpc server
	conn, err := grpc.Dial("192.168.1.100:5921", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	client := rpc.NewProcessClient(conn)

	req := rpc.StartRequest{
		Command: "php process.php",
	}
	resp, err := client.StartProcess(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}
