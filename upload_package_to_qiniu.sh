#!/usr/bin/env bash
# set -x -u
# 上传二进制包到七牛

# 打包
for i in linux darwin windows
do
    ./build.sh -p $i
    if [[ ! $? ]];then
        break
    fi
done

# 上传
for i in `ls gocron*.gz gocron*.zip`
do
    # 身份认证 qrsctl login <AccessKey> <SecretKey>
    # 上传文件 qrsctl put bucket  key srcFile
    qrsctl put github "gocron/${i}" $i
    if [[ ! $? ]];then
        break
    fi
    rm $i
done

echo '打包并上传成功'