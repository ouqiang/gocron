package ssh

import (
    "golang.org/x/crypto/ssh"
    "fmt"
    "net"
    "time"
    "errors"
)

type SSHConfig struct  {
    User string
    Password string
    Host string
    Port int
    ExecTimeout int// 执行超时时间
}

type Result struct {
    Output string
    Err error
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
        cmd += fmt.Sprintf(" & { sleep %d; eval 'kill  $!' &> /dev/null; }", sshConfig.ExecTimeout)
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
    config := &ssh.ClientConfig{
        User: sshConfig.User,
        Auth: []ssh.AuthMethod{
            ssh.Password(sshConfig.Password),
        },
        Timeout: 10 * time.Second,
        HostKeyCallback:func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
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