package utils

import (
    "crypto/md5"
    "encoding/hex"
    "math/rand"
    "os/exec"
    "time"
    "runtime"
    "github.com/Tang-RoseChild/mahonia"
    "strings"
)


// 执行shell命令，可设置执行超时时间
func ExecShellWithTimeout(timeout int, command string, args... string) (string, error)  {
    cmd := exec.Command(command, args...)

    // 不限制超时时间
    if timeout <= 0 {
        output ,err := cmd.CombinedOutput()
        return string(output), err
    }


    d := time.Duration(timeout) * time.Second
    timer := time.AfterFunc(d, func() {
        // 超时kill进程
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

// 判断当前系统是否是windows
func IsWindows() bool {
    return runtime.GOOS == "windows"
}

// GBK编码转换为UTF8
func GBK2UTF8(s string) (string, bool) {
    dec := mahonia.NewDecoder("gbk")

    return dec.ConvertStringOK(s)
}

// 转义json特殊字符
func EscapeJson(s string) string  {
    specialChars := []string{"\\", "\b","\f", "\n", "\r", "\t", "\"",}
    replaceChars := []string{ "\\\\", "\\b", "\\f", "\\n", "\\r", "\\t", "\\\"",}
    for i, v := range specialChars {
        s = strings.Replace(s, v, replaceChars[i], 1000)
    }

    return s
}