FROM alpine:3.7

ENV GOCRON_AGENT_VERSION=v1.5

RUN apk add --no-cache ca-certificates  tzdata bash \
    &&  mkdir -p /app \
    &&  wget -P /tmp  https://github.com/ouqiang/gocron/releases/download/${GOCRON_AGENT_VERSION}/gocron-node-${GOCRON_AGENT_VERSION}-linux-amd64.tar.gz \
    &&  cd /tmp \
    &&  tar  zvxf gocron-node-${GOCRON_AGENT_VERSION}-linux-amd64.tar.gz  \
    &&  mv /tmp/gocron-node-linux-amd64/gocron-node /app \
    &&  rm  -rf /tmp/* \
    &&  cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app
EXPOSE 5921

ENTRYPOINT ["/app/gocron-node", "-allow-root"]
