package utils

import (
	"os/exec"
	"math/rand"
	"time"
	"crypto/md5"
	"encoding/hex"
	"log"
)

// 执行shell命令
func ExecShell(command string, args... string) (string, error) {
	result, err := exec.Command(command, args...).CombinedOutput()

	return string(result), err
}

// 生成长度为length的随机字符串
func RandString(length int64) string {
	sources := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sourceLength := len(sources)
	var i int64 = 0
	for ; i < length; i++ {
		result = append(result, sources[r.Intn(sourceLength)])
	}

	return string(result)
}

// 生成32位MD5摘要
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))

	return hex.EncodeToString(m.Sum(nil))
}

// 生成0-max之间随机数
func RandNumber(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Intn(max)
}

// 日志记录
// todo 保存到哪里 文件,数据库还是elasticsearch?，暂时输出到终端
func RecordLog(v... interface{})  {
	log.Println(v)
}