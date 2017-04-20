package ssh

import (
    "golang.org/x/crypto/ssh"
    "fmt"
    "net"
    "time"
    "errors"
)


type HostAuthType int8  // 认证方式

const (
    HostPassword = 1   // 密码认证
    HostPublicKey = 2  // 公钥认证
)

const SSHConnectTimeout = 10


type SSHConfig struct  {
    AuthType HostAuthType
    User string
    Password string
    PrivateKey string
    Host string
    Port int
    ExecTimeout int// 执行超时时间
}

type Result struct {
    Output string
    Err error
}

func parseSSHConfig(sshConfig SSHConfig) (config *ssh.ClientConfig, err error) {
    timeout := SSHConnectTimeout * time.Second
    // 密码认证
    if sshConfig.AuthType == HostPassword {
        config = &ssh.ClientConfig{
            User: sshConfig.User,
            Auth: []ssh.AuthMethod{
                ssh.Password(sshConfig.Password),
            },
            Timeout: timeout,
            HostKeyCallback:func(hostname string, remote net.Addr, key ssh.PublicKey) error {
                return nil
            },
        }

        return
    }

    signer, err := ssh.ParsePrivateKey([]byte(sshConfig.PrivateKey))
    if err != nil {
        return
    }

    // 公钥认证
    config = &ssh.ClientConfig{
        User: sshConfig.User,
        Auth: []ssh.AuthMethod{
            ssh.PublicKeys(signer),
        },
        Timeout: timeout,
        HostKeyCallback:func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
    }

    return
}


// 执行shell命令
func Exec(sshConfig SSHConfig, cmd string) (output string, err error) {
    client, err := getClient(sshConfig)
    if err != nil {
        return "", err
    }
    defer client.Close()

    session, err := client.NewSession()

    if err != nil {
        return  "", err
    }
    defer session.Close()

    var resultChan chan Result = make(chan Result)
    var timeoutChan chan bool = make(chan bool)
    go func() {
        output, err := session.CombinedOutput(cmd)
        resultChan <- Result{string(output), err}
    }()
    go triggerTimeout(timeoutChan, sshConfig.ExecTimeout)
    select {
        case result := <- resultChan:
            output = result.Output
            err = result.Err
        case <- timeoutChan:
            output = ""
            err = errors.New("timeout")
    }

    return
}

func getClient(sshConfig SSHConfig) (*ssh.Client, error)  {
    config, err := parseSSHConfig(sshConfig)
    if err != nil {
        return nil, err
    }
    addr := fmt.Sprintf("%s:%d", sshConfig.Host, sshConfig.Port)

    return ssh.Dial("tcp", addr, config)
}


func triggerTimeout(ch chan bool, timeout int){
    // 最长执行时间不能超过24小时
    if timeout <= 0 || timeout > 86400 {
        timeout = 86400
    }
    time.Sleep(time.Duration(timeout) * time.Second)
    close(ch)
}

