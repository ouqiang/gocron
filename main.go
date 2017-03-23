package main

/*--------------------------------------------------------
 					定时任务调度
 	兼容Linux crontab时间格式语法，最小粒度可精确到每秒
 	支持通过HTTP、SSH协议触发任务执行
--------------------------------------------------------*/

import (
	"github.com/urfave/cli"
	"os"

	"github.com/ouqiang/cron-scheduler/cmd"
)

const AppVersion = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "scheduler"
	app.Usage = "schedule cron service"
	app.Version = AppVersion
	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
