package utils

import (
	"os/exec"
	"math/rand"
	"time"
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"runtime"
	"scheduler/utils/app"
)

// 检测环境
func CheckEnv()  {
	// ansible不支持安装在windows上, windows只能作为被控机
	if runtime.GOOS == "windows" {
		panic("不支持在windows上运行")
	}
	_, err := ExecShell("ansible", "--version")
	if err != nil {
		panic(err)
	}
	_, err = ExecShell("ansible-playbook", "--version")
	if err != nil {
		panic("ansible-playbook not found")
	}
}

// 创建安装锁文件
func CreateInstallLock()  {
	_, err := os.Create(app.ConfDir + "/install.lock")
	if err != nil {
		RecordLog("创建安装锁文件失败")
	}
}

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