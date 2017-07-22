// +build node
// 任务节点

package main

import (
	"github.com/ouqiang/gocron/modules/rpc/server"
	"os"
)

func main()  {
	var addr string
	if (len(os.Args) < 2) {
		addr = "0.0.0.0:5921"
	} else {
        addr = os.Args[1]
    }
	server.Start(addr)
}