package ansible

import (
    "bytes"
    "github.com/ouqiang/cron-scheduler/models"
    "github.com/ouqiang/cron-scheduler/modules/utils"
    "io/ioutil"
    "strconv"
    "sync"
)

// 主机名
var DefaultHosts *Hosts

type Hosts struct {
    sync.RWMutex
    filename string
}

func NewHosts(hostFilename string) *Hosts {
    h := &Hosts{sync.RWMutex{}, hostFilename}
    h.Write()

    return h
}

// 获取hosts文件名
func (h *Hosts) GetFilename() string {
    h.RLock()
    defer h.RUnlock()

    return h.filename
}

// 写入hosts
func (h *Hosts) Write() {
    host := new(models.Host)
    hostModels, err := host.List()
    if err != nil {
        utils.RecordLog(err)
        return
    }
    if len(hostModels) == 0 {
        utils.RecordLog("hosts内容为空")
        return
    }
    buffer := bytes.Buffer{}
    for _, hostModel := range hostModels {
        buffer.WriteString(strconv.Itoa(int(hostModel.Id)))
        buffer.WriteString(" ansible_ssh_host=")
        buffer.WriteString(hostModel.Name)
        buffer.WriteString(" ansible_ssh_port=")
        buffer.WriteString(strconv.Itoa(hostModel.Port))
        buffer.WriteString(" ansible_ssh_user=")
        buffer.WriteString(hostModel.Username)
        if hostModel.LoginType != models.PublicKey && hostModel.Password != "" {
            buffer.WriteString(" ansible_ssh_pass=")
            buffer.WriteString(hostModel.Password)
        }
        buffer.WriteString("\n")
    }
    h.Lock()
    defer h.Unlock()
    err = ioutil.WriteFile(h.filename, buffer.Bytes(), 0644)

    return
}
