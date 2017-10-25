FROM golang:1.9.1

MAINTAINER XUFEI <1842070912@qq.com>

# 配置时区
ENV TZ=Asia/Shanghai

WORKDIR /go/src/github.com/ouqiang/

RUN git clone https://github.com/ouqiang/gocron.git gocron

WORKDIR /go/src/github.com/ouqiang/gocron

RUN  go build -tags gocron -o gocron && go build -tags node -o gocron-node && \
    mkdir -p /go/src/github.com/ouqiang/gocron/conf

## add user
RUN useradd  gonode  -M  -s /sbin/nologin 

## add supervisor
RUN apt-get update && apt-get -y install supervisor vim &&  mkdir -p /var/log/supervisor 

## cp supervisord.conf

COPY supervisord.conf /etc/supervisord.conf

