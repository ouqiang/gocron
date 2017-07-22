#!/usr/bin/env bash

# set -x -u
# 任务节点打包, 生成压缩包 gocron-node.zip或gocron-node.tar.gz
# ./build-node.sh -p windows -a amd64
# 参数含义
# -p 指定平台(windows|linux|darwin)
# -a 指定体系架构(amd64|386), 默认amd64


# 目标平台 windows,linux,darwin
OS=''
# 目标平台架构
ARCH=''
# 应用名称
APP_NAME='gocron-node'
# 可执行文件名
EXEC_NAME=''
# 压缩包名称
COMPRESS_FILE=''


# -p 平台 -a 架构
while getopts "p:a:" OPT;
do
    case $OPT in
        p) OS=$OPTARG
        ;;
        a) ARCH=$OPTARG
        ;;
    esac
done

if [[ -z  $OS ]];then
    echo "平台不能为空"
    exit 1
fi

if [[ $OS != 'windows' && $OS != 'linux' && $OS != 'darwin' ]];then
    echo '平台错误，支持的平台 windows linux darmin(osx)'
    exit 1
fi

if [[ -z $ARCH ]];then
    ARCH='amd64'
fi

if [[ $ARCH != '386' && $ARCH != 'amd64' ]];then
    echo 'arch错误，仅支持 386 amd64'
    exit 1
fi

if [[ $OS = 'windows' ]];then
    EXEC_NAME=${APP_NAME}.exe
    COMPRESS_FILE=${APP_NAME}-${OS}-${ARCH}.zip
else
    EXEC_NAME=${APP_NAME}
    COMPRESS_FILE=${APP_NAME}-${OS}-${ARCH}.tar.gz
fi

echo '开始编译任务节点'
if [[ $OS = 'windows' ]];then
    GOOS=$OS GOARCH=$ARCH go build -tags node -ldflags '-w' -o $EXEC_NAME
else
    GOOS=$OS GOARCH=$ARCH go build -tags node -ldflags '-w' -o $EXEC_NAME
fi

if [[ $? != 0 ]];then
    exit 1
fi
echo '编译完成'

if [[ $OS = 'windows' ]];then
    zip -rq $COMPRESS_FILE  $EXEC_NAME
else
    tar czf $COMPRESS_FILE  $EXEC_NAME
fi


rm $EXEC_NAME

echo '打包完成'
echo '生成压缩文件--' $COMPRESS_FILE