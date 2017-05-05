[![Build Status](https://travis-ci.org/ouqiang/gocron.png)](https://travis-ci.org/ouqiang/gocron)
# gocron - 定时任务web管理系统

# 项目简介
使用Go语言开发的定时任务集中调度和管理系统, 用于替代Linux-crontab [查看文档](https://github.com/ouqiang/gocron/wiki)

## 功能特性
* 支持任务CURD
* crontab时间表达式，可精确到每秒
* 任务执行失败重试设置
* 任务超时设置
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
> Windows、Linux、OSX

### 环境要求
>  MySQL


## 安装
    
###  二进制安装
1. [下载](https://pan.baidu.com/s/1pLNRndx#list/path=%2Fgocron), 解压压缩包
2. `cd 解压目录`   
3. 启动  
    * Windows:  `gocron.exe web`            
    * Linux、OSX:  `./gocron web`
4. 浏览器访问 http://localhost:5920
### 源码安装
1. `go`语言版本1.7+
2. `go get -d https://github.com/ouqiang/gocron`
3. 编译 `go build`
4. 启动、访问方式同上
5. 生成压缩包（Windows: gocron.zip, 其他平台: gocron.tar.gz） 
> ./build.sh -p 平台 -a CPU架构 例 ./build.sh -p darwin -a amd64

### 启动可选参数

* -p 端口, 指定端口, 默认5920
* -e 指定运行环境, dev|test|prod, dev模式下可查看更多日志信息, 默认prod
* -h 查看帮助

## 安全
* 使用`https`访问保证数据传输安全, 可在web服务器如nginx中配置https，通过反向代理，访问内部的gocron
* 网站访问设置IP白名单
* SSH登录设置IP白名单
