package ansible

// ansible ad-hoc playbook命令封装

import (
	"os"
	"scheduler/utils"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"scheduler/utils/app"
)

func init()  {
		// ansible配置文件目录
		os.Setenv("ANSIBLE_CONFIG", app.ConfDir)
}

type Handler map[string]interface{}

type Playbook struct {
	Name string
	Hosts string
	Tasks []Handler
	Handlers []Handler
}

func(playbook *Playbook) SetHosts(hosts string) {
	playbook.Hosts = hosts
}

func(playbook *Playbook) SetName(name string) {
	playbook.Name = name
}

func(playbook *Playbook) AddTask(handler Handler) {
	playbook.Tasks = append(playbook.Tasks, handler)
}


func(playbook *Playbook) AddHandler(handler Handler) {
	playbook.Handlers = append(playbook.Handlers, handler)
}


/**
 * 执行ad-hoc
 * hosts  主机文件路径
 * module 调用模块
 * args   传递给模块的参数
*/
func ExecCommand(hostPath string, module string, args... string) (output string, err error) {
	if hostPath == "" || module == ""  {
		err = errors.New("参数不完整")
		return
	}
	commandArgs := []string{"-i", , hostPath, "-m", module}
	if len(args) != 0 {
		commandArgs = append(commandArgs, "-a", args...)
	}
	output, err = utils.ExecShell("ansible", commandArgs...)

	return
}

// 执行playbook
func ExecPlaybook(hostPath string, playbook Playbook) (result string, err error)  {
	data, err := yaml.Marshal([]Playbook{playbook})
	if err != nil {
		return
	}

	tmpFile, err := ioutil.TempFile(getTmpDir(), "playbook")
	if err != nil {
		return
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()
	_, err = tmpFile.Write(data)
	if err != nil {
		return
	}
	commandArgs := []string{"-i", hostPath, tmpFile.Name()}
	result, err = utils.ExecShell("ansible-playbook", commandArgs...)

	return
}

// 判断 获取临时目录，默认/dev/shm
func getTmpDir() string {
	dir := "/dev/shm"
	_, err := os.Stat(dir)
	if os.IsPermission(err) {
		return ""
	}

	return dir
}
