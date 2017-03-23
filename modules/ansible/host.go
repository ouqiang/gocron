package ansible

import (
	"scheduler/models"
	"sync"
	"io/ioutil"
	"bytes"
	"strconv"
)

// 主机名
var DefaultHosts *Hosts

type Hosts struct {
	sync.RWMutex
	hosts []models.Host
}

func(h *Hosts) Get() []models.Host {
	h.RLock()
	defer h.RUnlock()

	return h.hosts
}

func(h *Hosts) Set(hostsModel []models.Host)  {
	h.Lock()
	defer h.Unlock()

	h.hosts = hostsModel
}

// 获取hosts文件名
func(h *Hosts) GetHostFile() (filename string ,err error) {
	buffer := bytes.Buffer{}
	for _, hostModel := range(h.hosts) {
		buffer.WriteString(strconv.Itoa(int(hostModel.Id)))
		buffer.WriteString(" ansible_ssh_host=")
		buffer.WriteString(hostModel.Name)
		buffer.WriteString(" ansible_ssh_port=")
		buffer.WriteString(strconv.Itoa(hostModel.Port))
		buffer.WriteString(" ansible_ssh_user=")
		buffer.WriteString(hostModel.Username)
		if (hostModel.LoginType != models.PublicKey && hostModel.Password != "") {
			buffer.WriteString(" ansible_ssh_pass=")
			buffer.WriteString(hostModel.Password)
		}
		buffer.WriteString("\n")
	}
	tmpFile, err := ioutil.TempFile(GetTmpDir(), "host")
	if err != nil {
		return
	}

	defer func() {
		tmpFile.Close()
	}()

	_, err = tmpFile.WriteString(buffer.String())
	if err == nil {
		filename = tmpFile.Name()
	}

	return
}


