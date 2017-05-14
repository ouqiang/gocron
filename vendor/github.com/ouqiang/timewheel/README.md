# timewheel
Golang实现的时间轮


![时间轮](https://raw.githubusercontent.com/ouqiang/timewheel/master/timewheel.jpg)

# 安装

```shell
go get -u github.com/ouqiang/timewheel
```

# 使用

```go
package main

import (
    "github.com/ouqiang/timewheel"
    "time"
)

func main()  {
    // tick刻度为1秒, 3600个槽
    tw := timewheel.New(1 * time.Second, 3600)
    tw.Start()
    tw.Add(5 * time.Second, func() {
        // do something
    })
    tw.Add(10 * time.Minute, func() {
        // do something
    })
    tw.Add(35 * time.Hour, func() {
        // do something
    })
    // 停止
    tw.Stop()
}
```

