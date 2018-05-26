#!/usr/bin/env bash
 
# 生成压缩包 xx.tar.gz或xx.zip
# 使用 ./package.sh -a amd664 -p linux -v v2.0.0
 
# 任何命令返回非0值退出
set -o errexit
# 使用未定义的变量退出
set -o nounset
# 管道中任一命令执行失败退出
set -o pipefail

eval $(go env)

# 二进制文件名
BINARY_NAME=''
# main函数所在文件
MAIN_FILE=""
 
# 提取git最新tag作为应用版本
VERSION=''
# 最新git commit id
GIT_COMMIT_ID=''
 
# 外部输入的系统
INPUT_OS=()
# 外部输入的架构
INPUT_ARCH=()
# 未指定OS，默认值
DEFAULT_OS=${GOHOSTOS}
# 未指定ARCH,默认值
DEFAULT_ARCH=${GOHOSTARCH}
# 支持的系统
SUPPORT_OS=(linux darwin windows)
# 支持的架构
SUPPORT_ARCH=(386 amd64)
 
# 编译参数
LDFLAGS=''
# 需要打包的文件
INCLUDE_FILE=()
# 打包文件生成目录
PACKAGE_DIR=''
# 编译文件生成目录
BUILD_DIR=''
 
# 获取git 最新tag name
git_latest_tag() {
    local COMMIT_ID=""
    local TAG_NAME=""
    COMMIT_ID=`git rev-list --tags --max-count=1`
    TAG_NAME=`git describe --tags "${COMMIT_ID}"`
 
    echo ${TAG_NAME}
}
 
# 获取git 最新commit id
git_latest_commit() {
    echo "$(git rev-parse --short HEAD)"
}
 
# 打印信息
print_message() {
    echo "$1"
}
 
# 打印信息后推出
print_message_and_exit() {
    if [[ -n $1 ]]; then
        print_message "$1"
    fi
    exit 1
}
 
# 设置系统、CPU架构
set_os_arch() {
    if [[ ${#INPUT_OS[@]} = 0 ]];then
        INPUT_OS=("${DEFAULT_OS}")
    fi
 
    if [[ ${#INPUT_ARCH[@]} = 0 ]];then
        INPUT_ARCH=("${DEFAULT_ARCH}")
    fi
 
    for OS in "${INPUT_OS[@]}"; do
        if [[  ! "${SUPPORT_OS[*]}" =~ ${OS} ]]; then
            print_message_and_exit "不支持的系统${OS}"
        fi
    done
 
    for ARCH in "${INPUT_ARCH[@]}";do
        if [[ ! "${SUPPORT_ARCH[*]}" =~ ${ARCH} ]]; then
            print_message_and_exit "不支持的CPU架构${ARCH}"
        fi
    done
}
 
# 初始化
init() {
    set_os_arch
 
    if [[ -z "${VERSION}" ]];then
        VERSION=`git_latest_tag`
    fi
    GIT_COMMIT_ID=`git_latest_commit`
    LDFLAGS="-w -X 'main.AppVersion=${VERSION}' -X 'main.BuildDate=`date '+%Y-%m-%d %H:%M:%S'`' -X 'main.GitCommit=${GIT_COMMIT_ID}'"
 
    PACKAGE_DIR=${BINARY_NAME}-package
    BUILD_DIR=${BINARY_NAME}-build
 
    if [[ -d ${BUILD_DIR} ]];then
        rm -rf ${BUILD_DIR}
    fi
    if [[ -d ${PACKAGE_DIR} ]];then
        rm -rf ${PACKAGE_DIR}
    fi
 
    mkdir -p ${BUILD_DIR}
    mkdir -p ${PACKAGE_DIR}
}
 
# 编译
build() {
    local FILENAME=''
    for OS in "${INPUT_OS[@]}";do
        for ARCH in "${INPUT_ARCH[@]}";do
            if [[ "${OS}" = "windows"  ]];then
                FILENAME=${BINARY_NAME}.exe
            else
                FILENAME=${BINARY_NAME}
            fi
            env CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME}-${OS}-${ARCH}/${FILENAME} ${MAIN_FILE}
        done
    done
}
 
# 打包
package_binary() {
    cd ${BUILD_DIR}
 
    for OS in "${INPUT_OS[@]}";do
        for ARCH in "${INPUT_ARCH[@]}";do
        package_file ${BINARY_NAME}-${OS}-${ARCH}
        if [[ "${OS}" = "windows" ]];then
            zip -rq ../${PACKAGE_DIR}/${BINARY_NAME}-${VERSION}-${OS}-${ARCH}.zip ${BINARY_NAME}-${OS}-${ARCH}
        else
            tar czf ../${PACKAGE_DIR}/${BINARY_NAME}-${VERSION}-${OS}-${ARCH}.tar.gz ${BINARY_NAME}-${OS}-${ARCH}
        fi
        done
    done
 
    cd ${OLDPWD}
}
 
# 打包文件
package_file() {
    if [[ "${#INCLUDE_FILE[@]}" = "0" ]];then
        return
    fi
    for item in "${INCLUDE_FILE[@]}"; do
            cp -r ../${item} $1
    done
}
 
# 清理
clean() {
    if [[ -d ${BUILD_DIR} ]];then
        rm -rf ${BUILD_DIR}
    fi
}
 
# 运行
run() {
    init
    build
    package_binary
    clean
}

package_gocron() {
    BINARY_NAME='gocron'
    MAIN_FILE="./cmd/gocron/gocron.go"
    INCLUDE_FILE=()


    run
}

package_gocron_node() {
    BINARY_NAME='gocron-node'
    MAIN_FILE="./cmd/node/node.go"
    INCLUDE_FILE=()

    run
}
 
# p 平台 linux darwin windows
# a 架构 386 amd64
# v 版本号  默认取git最新tag
while getopts "p:a:v:" OPT;
do
    case ${OPT} in
    p) IPS=',' read -r -a INPUT_OS <<< "${OPTARG}"
    ;;
    a) IPS=',' read -r -a INPUT_ARCH <<< "${OPTARG}"
    ;;
    v) VERSION=$OPTARG
    ;;
    *)
    ;;
    esac
done
 
package_gocron
package_gocron_node

