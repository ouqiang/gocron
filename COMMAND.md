##### 常用命令记录

生成grpc文件
```shell
# 安装依赖
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# 生成文件
$ protoc --go_out=internal/modules/rpc/proto/ --go-grpc_out=internal/modules/rpc/proto/ --go-grpc_opt=require_unimplemented_servers=false internal/modules/rpc/proto/process.proto
```

statik.go文件生成
1. 先下载并安装 statik 这个工具
> go get -d github.com/rakyll/statik

> go install github.com/rakyll/statik

注意将 $GOPATH/bin 加入到 PATH 环境变量中。

2. 编译vue文件
> cd web/vue && npm run build

3. 生成statik.go
> go generate cmd/gocron/gocron.go