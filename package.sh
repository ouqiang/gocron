#!/usr/bin/env bash

#set -x -u
# 构建应用, 生成压缩包 gocron.zip或gocron.tar.gz
# ./package.sh -v 1.4

VERSION=''
GOCRON_APP_NAME='gocron'
GOCRON_NODE_APP_NAME='gocron-node'
PACKAGE_DIR='./packages'


# 用法
usage() {
    echo 'usage: ./package.sh -v version'
}

# 初始化
init() {
    rm -rf ${PACKAGE_DIR}
    mkdir -p ${PACKAGE_DIR}
}

# 构建应用
build() {
    make -f makefile.cross-compiles VERSION=${VERSION}

    if [[ $? -ne 0 ]];then
        echo 'make error'
        exit 1
    fi
}

# 清理
clean() {
    make -f makefile.cross-compiles clean
}

# 打包gocron
package_gocron() {
    local OS=$1
    local GOCRON_COMPRESS_FILE=''
    local PLATFORM_NAME=${GOCRON_APP_NAME}_${OS}_amd64

    if [[ ! -d ${PACKAGE_DIR}/${PLATFORM_NAME} ]];then
        mkdir ${PACKAGE_DIR}/${PLATFORM_NAME}
    fi

    for file in public templates LICENSE README.md Dockerfile-release; do
        cp -r ${file} ${PLATFORM_NAME}
    done

    if [[ ${OS} = 'windows' ]];then
        GOCRON_COMPRESS_FILE=${GOCRON_APP_NAME}-v${VERSION}-${OS}-amd64.zip
        zip -rq ${PACKAGE_DIR}/${PLATFORM_NAME}/${GOCRON_COMPRESS_FILE} ${PLATFORM_NAME}
    else
        GOCRON_COMPRESS_FILE=${GOCRON_APP_NAME}-v${VERSION}-${OS}-amd64.tar.gz
        tar czf ${PACKAGE_DIR}/${PLATFORM_NAME}/${GOCRON_COMPRESS_FILE} ${PLATFORM_NAME}
    fi
}

# 打包gocron-node
package_gocron_node() {
    local OS=$1
    local GOCRON_NODE_COMPRESS_FILE=''
    local PLATFORM_NAME=${GOCRON_NODE_APP_NAME}_${OS}_amd64

    if [[ ! -d ${PACKAGE_DIR}/${PLATFORM_NAME} ]];then
        mkdir ${PACKAGE_DIR}/${PLATFORM_NAME}
    fi

    if [[ ${OS} = 'windows' ]];then
        GOCRON_NODE_COMPRESS_FILE=${GOCRON_NODE_APP_NAME}-v${VERSION}-${OS}-amd64.zip
        zip -rq ${PACKAGE_DIR}/${PLATFORM_NAME}/${GOCRON_NODE_COMPRESS_FILE} ${PLATFORM_NAME}
    else
        GOCRON_NODE_COMPRESS_FILE=${GOCRON_NODE_APP_NAME}-v${VERSION}-${OS}-amd64.tar.gz
        tar czf ${PACKAGE_DIR}/${PLATFORM_NAME}/${GOCRON_NODE_COMPRESS_FILE} ${PLATFORM_NAME}
    fi
}


package_multi() {
    for os in darwin linux windows; do
        package_gocron ${os}
        package_gocron_node ${os}
    done
}

while getopts "v:" OPT; do
    case ${OPT} in
        v) VERSION=${OPTARG}
        ;;
    esac
done

if [[ -z ${VERSION} ]]; then
    usage
    exit 1
fi


run() {
    init
    build
    package_multi
    clean
}


run




