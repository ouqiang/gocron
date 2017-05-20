package cmd

import (
    "github.com/urfave/cli"
    "fmt"
    "github.com/ouqiang/gocron/modules/utils"
    "github.com/ouqiang/gocron/modules/app"
    "os"
    "syscall"
)

var CmdServ = cli.Command{
    Name:   "serv",
    Usage:  "manage gocron, ./gocron serv -s stop|status",
    Action: runServ,
    Flags: []cli.Flag{
        cli.StringFlag{
            Name: "s",
            Value:"",
            Usage: "stop|status",
        },
    },
}

func runServ(ctx *cli.Context)  {
    if utils.IsWindows() {
        fmt.Println("not support on windows")
        return
    }
    option := ctx.String("s")
    if !utils.InStringSlice([]string{"stop", "status"}, option) {
        fmt.Println("invalid option")
        return
    }
    app.InitEnv()
    pid := app.GetPid()
    if pid <= 0 {
        fmt.Println("not running")
        return
    }
    process ,err := os.FindProcess(pid)
    if err != nil {
        fmt.Println("not running", err)
        return
    }
    switch option {
        case "stop":
            stop(process)
        case "status":
            status(process)
    }
}

func stop(process *os.Process)  {
    fmt.Println("stopping gocron......")
    err := process.Signal(syscall.SIGTERM)
    if err != nil {
        fmt.Println("failed to kill process", err)
    } else {
        fmt.Println("stopped")
    }
}


func status(process *os.Process)  {
    fmt.Printf("running, pid-[%d]\n", process.Pid)
}