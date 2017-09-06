// +build node
// 任务节点

package main

import (
	"github.com/ouqiang/gocron/modules/rpc/server"
    "flag"
    "runtime"
    "os"
    "fmt"
    "strings"
)

const AppVersion = "1.1"

func main()  {
	var serverAddr string
    var allowRoot bool
    var version bool
    var keyFile string
    var certFile string
    flag.BoolVar(&allowRoot, "allow-root", false, "./gocron-node -allow-root")
    flag.StringVar(&serverAddr, "s", "0.0.0.0:5921", "./gocron-node -s ip:port")
    flag.StringVar(&certFile, "cert-file", "", "./gocron-node -cert-file path")
    flag.StringVar(&keyFile, "key-file", "", "./gocron-node -key-file path")
    flag.BoolVar(&version, "v", false, "./gocron-node -v")
    flag.Parse()

    if version {
        fmt.Println(AppVersion)
        os.Exit(0)
    }

    certFile = strings.TrimSpace(certFile)
    keyFile = strings.TrimSpace(keyFile)

    if certFile != "" && keyFile == "" {
        fmt.Println("missing argument key-file")
        return
    }

    if keyFile != "" && certFile == "" {
        fmt.Println("missing argument cert-file")
        return
    }

    if runtime.GOOS != "windows" && os.Getuid() == 0 && !allowRoot   {
        fmt.Println("Do not run gocron-node as root user")
        os.Exit(1)
    }



	server.Start(serverAddr, certFile, keyFile)
}