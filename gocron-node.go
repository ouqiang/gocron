// +build node
// 任务节点

package main

import (
	"github.com/ouqiang/gocron/modules/rpc/server"
    "flag"
    "runtime"
    "os"
    "fmt"
)

const AppVersion = "1.2"

func main()  {
	var serverAddr string
    var allowRoot bool
    var version bool
    flag.BoolVar(&allowRoot, "allow-root", false, "./gocron-node -allow-root")
    flag.StringVar(&serverAddr, "s", "0.0.0.0:5921", "./gocron-node -s ip:port")
    flag.BoolVar(&version, "v", false, "./gocron-node -v")
    flag.Parse()

    if version {
        fmt.Println(AppVersion)
        os.Exit(0)
    }

    if runtime.GOOS != "windows" && os.Getuid() == 0 && !allowRoot   {
        fmt.Println("Do not run gocron-node as root user")
        os.Exit(1)
    }


	server.Start(serverAddr)
}