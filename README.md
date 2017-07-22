[![Build Status](https://travis-ci.org/ouqiang/gocron.png)](https://travis-ci.org/ouqiang/gocron)
# gocron - 定时任务管理系统

# 项目简介
使用Go语言开发的定时任务集中调度和管理系统, 用于替代Linux-crontab [查看文档](https://github.com/ouqiang/gocron/wiki)  

原有的延时任务拆分为独立项目[延迟队列](https://github.com/ouqiang/delay-queue)  

## 功能特性
* Web界面管理定时任务, 支持动态添加、删除任务
* crontab时间表达式，精确到秒
* 任务执行失败重试设置
* 任务超时设置
* 任务依赖配置
* 任务类型
    * shell任务
    > 在任务节点上执行shell命令
    * HTTP任务
    > 访问指定的URL地址, 由调度器直接执行, 不依赖任务节点
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
[v1.0](https://github.com/ouqiang/gocron/releases/tag/v1.0)


## 安装

###  二进制安装
1. 解压压缩包   
2. `cd 解压目录`   
3. 启动        
* 调度器启动        
  * Windows: `gocron.exe web`   
  * Linux、Mac OS:  `./gocron web`
* 任务节点启动
  * Windows:  `gocron-node.exe ip:port (默认0.0.0.0:5921)`            
  * Linux、Mac OS:  `./gocron-node ip:port (默认0.0.0.0:5921)`   
4. 浏览器访问 http://localhost:5920

### 源码安装
1. `go`语言版本1.7+
2. `go get -d github.com/ouqiang/gocron`
3. 编译 
    * 调度器 `go build -tags gocron -o gocron`
    * 任务节点 `go build -tags node -o gocron-node`
4. 启动、访问方式同上

### 命令

* gocron web
    * --host 默认0.0.0.0
    * -p 端口, 指定端口, 默认5920
    * -e 指定运行环境, dev|test|prod, dev模式下可查看更多日志信息, 默认prod
    * -d 后台运行
    * -h 查看帮助
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
