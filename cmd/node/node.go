// Command gocron-node
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ouqiang/gocron/internal/modules/rpc/auth"
	"github.com/ouqiang/gocron/internal/modules/rpc/server"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"github.com/ouqiang/goutil"
)

var (
	AppVersion, BuildDate, GitCommit string
)

func main() {
	var serverAddr string
	var allowRoot bool
	var version bool
	var CAFile string
	var certFile string
	var keyFile string
	var enableTLS bool
	flag.BoolVar(&allowRoot, "allow-root", false, "./gocron-node -allow-root")
	flag.StringVar(&serverAddr, "s", "0.0.0.0:5921", "./gocron-node -s ip:port")
	flag.BoolVar(&version, "v", false, "./gocron-node -v")
	flag.BoolVar(&enableTLS, "enable-tls", false, "./gocron-node -enable-tls")
	flag.StringVar(&CAFile, "ca-file", "", "./gocron-node -ca-file path")
	flag.StringVar(&certFile, "cert-file", "", "./gocron-node -cert-file path")
	flag.StringVar(&keyFile, "key-file", "", "./gocron-node -key-file path")
	flag.Parse()

	if version {
		goutil.PrintAppVersion(AppVersion, GitCommit, BuildDate)
		return
	}

	if enableTLS {
		if !utils.FileExist(CAFile) {
			fmt.Printf("failed to read ca cert file: %s", CAFile)
			return
		}
		if !utils.FileExist(certFile) {
			fmt.Printf("failed to read server cert file: %s", certFile)
			return
		}
		if !utils.FileExist(keyFile) {
			fmt.Printf("failed to read server key file: %s", keyFile)
			return
		}
	}

	certificate := auth.Certificate{
		CAFile:   strings.TrimSpace(CAFile),
		CertFile: strings.TrimSpace(certFile),
		KeyFile:  strings.TrimSpace(keyFile),
	}

	if runtime.GOOS != "windows" && os.Getuid() == 0 && !allowRoot {
		fmt.Println("Do not run gocron-node as root user")
		return
	}

	server.Start(serverAddr, enableTLS, certificate)
}
