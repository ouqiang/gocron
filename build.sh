#!/usr/bin/env bash

# set -x -u
# 构建应用, 生成压缩包 gocron.zip或gocron.tar.gz
# ./build.sh -p windows -a amd64 -v 1.4
# 参数含义
# -p 指定平台(windows|linux|darwin)
# -a 指定体系架构(amd64|386), 默认amd64
# -v 版本号

TEMP_DIR=`date +%s`-temp-`echo $RANDOM`

# 目标平台 windows,linux,darwin
OS=''
# 目标平台架构
ARCH=''
# 应用名称
APP_NAME='gocron'
# 版本号
VERSION=''
# 可执行文件名
EXEC_NAME=''
# 压缩包名称
COMPRESS_FILE=''


# -p 平台 -a 架构
while getopts "p:a:v:" OPT;
do
    case $OPT in
        p) OS=$OPTARG
        ;;
        a) ARCH=$OPTARG
        ;;
        v) VERSION=$OPTARG
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

if [[ -z $VERSION ]];then
    echo '版本号不能为空'
    exit 1
fi

echo '开始编译调度器'
if [[ $OS = 'windows' ]];then
    GOOS=$OS GOARCH=$ARCH go build -tags gocron -ldflags '-w'
else
    GOOS=$OS GOARCH=$ARCH go build -tags gocron -ldflags '-w'
fi

if [[ $? != 0 ]];then
    exit 1
fi
echo '编译完成'

if [[ $OS = 'windows' ]];then
    EXEC_NAME=${APP_NAME}.exe
    COMPRESS_FILE=${APP_NAME}-v${VERSION}-${OS}-${ARCH}.zip
else
    EXEC_NAME=${APP_NAME}
    COMPRESS_FILE=${APP_NAME}-v${VERSION}-${OS}-${ARCH}.tar.gz
fi

mkdir -p $TEMP_DIR/$APP_NAME
if [[ $? != 0 ]]; then
    exit 1
fi

# 需要打包的文件
PACKAGE_FILENAME=(public templates ${EXEC_NAME})

echo '复制文件到临时目录'
# 复制文件到临时目录
for i in ${PACKAGE_FILENAME[*]}
do
    cp -r $i $TEMP_DIR/$APP_NAME
done

echo '压缩文件'
# 压缩文件
cd $TEMP_DIR
if [[ $OS = 'windows' ]];then
    zip -rq $COMPRESS_FILE *
else
    tar czf $COMPRESS_FILE *
fi
mv $COMPRESS_FILE ../
cd ../

rm $EXEC_NAME
rm -rf $TEMP_DIR

echo '打包完成'
echo '生成压缩文件--' $COMPRESS_FILE