// +build node
// 任务节点

package main

import (
	"github.com/ouqiang/gocron/modules/rpc/server"
	"os"
	"fmt"
)

func main()  {
	var addr string
	if (len(os.Args) < 2) {
		fmt.Println("usage ./gocron-node addr:port")
		os.Exit(1)
	}
	addr = os.Args[1]
	server.Start(addr)
}