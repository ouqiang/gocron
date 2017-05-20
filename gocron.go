package main

/*--------------------------------------------------------
                定时任务调度
 	兼容Linux crontab时间格式语法，最小粒度可精确到每秒
 	支持通过HTTP、SSH协议执行任务
--------------------------------------------------------*/

import (
    "github.com/urfave/cli"
    "os"

    "github.com/ouqiang/gocron/cmd"
)

const AppVersion = "0.3"

func main() {
    app := cli.NewApp()
    app.Name = "gocron"
    app.Usage = "gocron service"
    app.Version = AppVersion
    app.Commands = []cli.Command{
        cmd.CmdWeb,
        cmd.CmdServ,
    }
    app.Flags = append(app.Flags, []cli.Flag{}...)
    app.Run(os.Args)
}
