// +build gocron
// 调度中心

package main

import (
	"github.com/urfave/cli"
	"os"

	"github.com/ouqiang/gocron/cmd"
)

const AppVersion = "1.2.2"

func main() {
	app := cli.NewApp()
	app.Name = "gocron"
	app.Usage = "gocron service"
	app.Version = AppVersion
	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
