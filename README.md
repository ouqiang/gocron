[![Build Status](https://travis-ci.org/ouqiang/gocron.png)](https://travis-ci.org/ouqiang/gocron)
# gocron - 定时任务web管理系统

# 项目简介
使用Go语言开发的定时任务集中调度和管理系统, 用于替代Linux-crontab [查看文档](https://github.com/ouqiang/gocron/wiki)

## 功能特性
* 支持任务CURD
* crontab时间表达式，精确到秒
* 任务执行失败重试设置
* 任务超时设置
* 延时任务
* 任务执行方式
    * 调用本机系统命令  
    * 通过SSH执行远程命令
    * 执行HTTP-GET请求
* 查看任务执行日志
* 任务执行结果通知, 支持邮件、Slack

### 截图
![任务](https://raw.githubusercontent.com/ouqiang/gocron/master/screenshot_task.png)
![Slack](https://raw.githubusercontent.com/ouqiang/gocron/master/screenshot_slack.png)
    
### 支持平台
> Windows、Linux、Mac OS

### 环境要求
>  MySQL


## 下载
* [Linux-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-linux-amd64.tar.gz)
* [Mac OS-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-darwin-amd64.tar.gz)
* [Windows-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-windows-amd64.zip)

## 安装

###  二进制安装
1. 解压压缩包    
2. `cd 解压目录`   
3. 启动  
    * Windows:  `gocron.exe web`            
    * Linux、Mac OS:  `./gocron web`
4. 浏览器访问 http://localhost:5920
### 源码安装
1. `go`语言版本1.7+
2. `go get -d github.com/ouqiang/gocron`
3. 编译 `go build`
4. 启动、访问方式同上

### 命令

* gocron web
    * -p 端口, 指定端口, 默认5920
    * -e 指定运行环境, dev|test|prod, dev模式下可查看更多日志信息, 默认prod
    * -d 后台运行
    * -h 查看帮助
* gocron serv 
    * -s stop|status stop:停止gocron status:查看运行状态
    

## 安全
* 使用`https`访问保证数据传输安全, 可在web服务器如nginx中配置https，通过反向代理，访问内部的gocron
* 网站访问设置IP白名单
* SSH登录设置IP白名单

## 程序使用的组件
* web框架 [Macaron](http://go-macaron.com/)
* 定时任务调度 [cron](https://github.com/robfig/cron)
* ORM [Xorm](https://github.com/go-xorm/xorm)
* UI框架 [Semantic UI](https://semantic-ui.com/)
* 依赖管理(所有依赖包放入vendor目录) [govendor](https://github.com/kardianos/govendor)

## 贡献
欢迎提交PR
