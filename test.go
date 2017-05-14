package main

import (
    "github.com/ouqiang/timewheel"
    "time"
    "fmt"
    "github.com/ouqiang/gocron/models"
)

func main()  {
    // tick刻度为1秒, 3600个槽
    tw := timewheel.New(1 * time.Second, 3600)
    tw.Start()
    t := time.Now()
    tw.Add(5 * time.Second, func() {
        fmt.Println("5分钟", time.Now(), t.Add(5 * time.Second).Format(models.DefaultTimeFormat))
    })
    tw.Add(10 * time.Second, func() {
        fmt.Println("10分钟", time.Now(), t.Add(10 * time.Second).Format(models.DefaultTimeFormat))
    })
    tw.Add(35 * time.Second, func() {
        fmt.Println("35分钟", time.Now(), t.Add(35 * time.Second).Format(models.DefaultTimeFormat))
    })
    tw.Add(178 * time.Second, func() {
        fmt.Println("178分钟", time.Now(), t.Add(178 * time.Second).Format(models.DefaultTimeFormat))
    })
    tw.Add(27 * time.Second, func() {
        fmt.Println("27分钟", time.Now(), t.Add(27 * time.Second).Format(models.DefaultTimeFormat))
    })
    tw.Add(78 * time.Second, func() {
        fmt.Println("78分钟", time.Now(), t.Add(78 * time.Second).Format(models.DefaultTimeFormat))
    })
    tw.Add(3 * time.Second, func() {
        fmt.Println("3分钟", time.Now(), t.Add(3 * time.Second).Format(models.DefaultTimeFormat))
    })
    select {}
}
