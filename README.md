[![Build Status](https://travis-ci.org/ouqiang/gocron.png)](https://travis-ci.org/ouqiang/gocron)
# gocron - 定时任务管理系统

# 项目简介
使用Go语言开发的定时任务集中调度和管理系统, 用于替代Linux-crontab [查看文档](https://github.com/ouqiang/gocron/wiki)

## 功能特性
* Web界面管理定时任务, 支持动态添加、删除、编辑任务
* crontab时间表达式，精确到秒
* 任务执行失败重试设置
* 任务超时设置
* 延时任务
* 任务依赖配置
* 任务类型
    * shell任务
    > 在远程服务器上执行shell命令
    * HTTP任务
    > 访问指定的URL地址
* 查看任务执行日志
* 任务执行结果通知, 支持邮件、Slack

### 截图
![流程图](https://raw.githubusercontent.com/ouqiang/gocron/master/scheduler.png)
![任务](https://raw.githubusercontent.com/ouqiang/gocron/master/screenshot_task.png)
![Slack](https://raw.githubusercontent.com/ouqiang/gocron/master/screenshot_slack.png)
    
### 支持平台
> Windows、Linux、Mac OS

### 环境要求
>  MySQL


## 下载
* 调度器(管理后台)
    * [Linux-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-linux-amd64.tar.gz)
    * [Mac OS-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-darwin-amd64.tar.gz)
    * [Windows-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-windows-amd64.zip)
* 任务执行器(安装在远程主机上, 执行shell命令需安装)
     * [Linux-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-node-linux-amd64.tar.gz)
     * [Mac OS-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-node-darwin-amd64.tar.gz)
     * [Windows-64位](http://opns468ov.bkt.clouddn.com/gocron/gocron-node-windows-amd64.zip)


## 安装

###  二进制安装
> Windows平台默认后台运行)
1. 解压压缩包   
2. `cd 解压目录`   
3. 启动        
    * 调度器启动      
        * Windows: `gocron.exe web`    
        * Linux、Mac OS:  `./gocron web`
    * 任务执行器启动
        * Windows:  `gocron-node.exe ip:port (默认0.0.0.0:5921)`            
        * Linux、Mac OS:  `./gocron-node ip:port (默认0.0.0.0:5921)`   
4. 浏览器访问 http://localhost:5920
### 源码安装
1. `go`语言版本1.7+
2. `go get -d github.com/ouqiang/gocron`
3. 编译 
    * 调度器 `go build -tags gocron -o gocron`
    * 任务执行器 `go build -tags node -o gocron-node`
4. 启动、访问方式同上

### 命令

* gocron web
    * --host 默认0.0.0.0
    * -p 端口, 指定端口, 默认5920
    * -e 指定运行环境, dev|test|prod, dev模式下可查看更多日志信息, 默认prod
    * -d 后台运行
    * -h 查看帮助
* gocron serv 
    * -s stop|status stop:停止gocron status:查看运行状态
* gocron-node ip:port, 默认0.0.0.0:5921 

## 程序使用的组件
* web框架 [Macaron](http://go-macaron.com/)
* 定时任务调度 [Cron](https://github.com/robfig/cron)
* ORM [Xorm](https://github.com/go-xorm/xorm)
* UI框架 [Semantic UI](https://semantic-ui.com/)
* 依赖管理(所有依赖包放入vendor目录) [Govendor](https://github.com/kardianos/govendor)
* RPC框架 [gRPC](https://github.com/grpc/grpc)

## 反馈
提交[issue](https://github.com/ouqiang/gocron/issues/new)
