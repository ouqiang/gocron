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
    > 在任务节点上执行shell命令, 支持任务同时在多个节点上运行
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
[v1.2.1](https://github.com/ouqiang/gocron/releases/tag/v1.2.1)  

[版本升级](https://github.com/ouqiang/gocron/wiki/版本升级)

## 安装

###  二进制安装
1. 解压压缩包   
2. `cd 解压目录`   
3. 启动        
* 调度器启动        
  * Windows: `gocron.exe web`   
  * Linux、Mac OS:  `./gocron web`
* 任务节点启动, 默认监听0.0.0.0:5921
  * Windows:  `gocron-node.exe`
  * Linux、Mac OS:  `./gocron-node`
4. 浏览器访问 http://localhost:5920

### 源码安装
1. `go`语言版本1.7+
2. `go get -d github.com/ouqiang/gocron`
3. 编译 
    * 调度器 `go build -tags gocron -o gocron`
    * 任务节点 `go build -tags node -o gocron-node`
4. 启动、访问方式同上

### 命令

* gocron
    * -v 查看版本

* gocron web
    * --host 默认0.0.0.0
    * -p 端口, 指定端口, 默认5920
    * -e 指定运行环境, dev|test|prod, dev模式下可查看更多日志信息, 默认prod
    * -h 查看帮助
* gocron-node
    * -allow-root *nix平台允许以root用户运行
    * -s ip:port 监听地址
    * -cert-file 证书文件
    * -key-file  私钥文件
    * -h 查看帮助
    * -v 查看版本

## To Do List
- [x] 版本升级
- [x] 批量开启、关闭、删除任务
- [x] 调度器与任务节点通信支持https
- [ ] 任务分组
- [ ] 多用户
- [ ] 权限控制
- [ ] 新增任务API接口

## 程序使用的组件
* Web框架 [Macaron](http://go-macaron.com/)
* 定时任务调度 [Cron](https://github.com/robfig/cron)
* ORM [Xorm](https://github.com/go-xorm/xorm)
* UI框架 [Semantic UI](https://semantic-ui.com/)
* 依赖管理(所有依赖包放入vendor目录) [Govendor](https://github.com/kardianos/govendor)
* RPC框架 [gRPC](https://github.com/grpc/grpc)

## 反馈
提交[issue](https://github.com/ouqiang/gocron/issues/new)

## ChangeLog

v1.2
--------
* 用户登录页增加图形验证码
* 支持从旧版本升级
* 任务批量开启、关闭、删除
* 调度器与任务节点支持HTTPS通信
* 修复任务列表页总记录数显示错误



v1.1
--------

* 任务可同时在多个节点上运行
* *nix平台默认禁止以root用户运行任务节点
* 子任务命令中增加预定义占位符, 子任务可根据主任务运行结果执行相应操作
* 删除守护进程模块
* Web访问日志输出到终端
