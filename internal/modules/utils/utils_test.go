package utils

import (
	"fmt"
	rpc "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"golang.org/x/net/context"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRandString(t *testing.T) {
	str := RandString(32)
	if len(str) != 32 {
		t.Fatalf("长度不匹配,目标长度32, 实际%d-%s", len(str), str)
	}
}

func TestMd5(t *testing.T) {
	str := Md5("123456")
	if len(str) != 32 {
		t.Fatalf("长度不匹配,目标长度32, 实际%d-%s", len(str), str)
	}
}

func TestRandNumber(t *testing.T) {
	num := RandNumber(10000)
	if num <= 0 || num >= 10000 {
		t.Fatalf("随机数不在有效范围内-%d", num)
	}
}

func TestStartWorker(t *testing.T) {
	req := rpc.StartRequest{
		Command: "php /var/www/go/gocron/process.php",
	}
	pid, err := StartWorker(context.Background(), &req)
	t.Log(pid, err)
}

func TestStopWorker(t *testing.T) {
	err := StopWorker(int64(15780))
	t.Log(err)
}

func TestBaseDir(t *testing.T) {
	absPath, _ := filepath.Abs("gocron/process")

	dir := path.Dir(absPath)
	dir2 := path.Dir("gocron/1/2/3/4/5/process")

	a1, _ := filepath.Abs(dir2)

	_, err := os.Stat(a1)
	if err != nil || os.IsNotExist(err) {
		_ = os.MkdirAll(a1, 0666)
	}
	t.Log(err)
	fmt.Println(dir, absPath, dir2, a1, os.IsNotExist(err))
}

func TestWorkerStateCheck(t *testing.T) {
	str, err := WorkerStateCheck(int64(0))
	t.Log("err: ", str, err)
}
