#!/usr/bin/env bash
# set -x -u
# 上传二进制包到七牛

if [[ -z $QINIU_ACCESS_KEY || -z  $QINIU_SECRET_KEY || -z $QINIU_URL ]];then
    echo 'QINIU_ACCESS_KEY | QINIU_SECRET_KEY | QINIU_URL is need'
    exit 1
fi

# 打包
for i in linux darwin windows
do
    ./build.sh -p $i
    if [[ $? != 0 ]];then
        break
    fi
done

# 身份认证
qrsctl login $QINIU_ACCESS_KEY $QINIU_SECRET_KEY

# 上传
for i in `ls gocron*.gz gocron*.zip`
do
    # 上传文件 qrsctl put bucket  key srcFile
    KEY=gocron/$i
    qrsctl put github $KEY $i
    if [[ $? != 0 ]];then
        break
    fi
    echo "刷新七牛CDN-" $QINIU_URL/$KEY
    qrsctl cdn/refresh $QINIU_URL/$KEY
    rm $i
done

echo '打包并上传成功'