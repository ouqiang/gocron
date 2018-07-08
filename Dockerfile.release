FROM alpine:3.7

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -g app app

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

COPY gocron .

RUN chown -R app:app ./

EXPOSE 5920

USER app

ENTRYPOINT ["/app/gocron", "web"]
