package utils

import (
    "crypto/md5"
    "encoding/hex"
    "math/rand"
    "os/exec"
    "time"
)

// 执行shell命令
func ExecShell(command string, args ...string) (string, error) {
    result, err := exec.Command(command, args...).CombinedOutput()

    return string(result), err
}

// 执行shell命令，可设置执行超时时间
func ExecShellWithTimeout(timeout int, command string, args... string) (string, error)  {
    cmd := exec.Command(command, args...)
    d := time.Duration(timeout) * time.Second
    timer := time.AfterFunc(d, func() {
        // 执行超时kill进程
        cmd.Process.Kill()
    })
    output ,err := cmd.CombinedOutput()
    timer.Stop()

    return string(output), err
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