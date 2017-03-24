package ansible

// ansible ad-hoc playbook命令封装

import (
	"os"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/ouqiang/cron-scheduler/modules/utils"
)

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
 * hosts  主机名 逗号分隔
 * module 调用模块
 * args   传递给模块的参数
*/
func ExecCommand(hosts string, module string, args... string) (output string, err error) {
	if hosts== "" || module == ""  {
		err = errors.New("参数不完整")
		return
	}
	hostFile := DefaultHosts.GetFilename()
	commandArgs := []string{hosts, "-i",  hostFile, "-m", module}
	if len(args) != 0 {
		commandArgs = append(commandArgs, "-a")
		commandArgs = append(commandArgs,  args...)
	}
	output, err = utils.ExecShell("ansible", commandArgs...)

	return
}

// 执行playbook
func ExecPlaybook(playbook Playbook) (result string, err error)  {
	data, err := yaml.Marshal([]Playbook{playbook})
	if err != nil {
		return
	}

	playbookFile, err := ioutil.TempFile(GetTmpDir(), "playbook")
	if err != nil {
		return
	}
	hostFile := DefaultHosts.GetFilename()
	defer func() {
		playbookFile.Close()
		os.Remove(playbookFile.Name())
	}()
	_, err = playbookFile.Write(data)
	if err != nil {
		return
	}
	commandArgs := []string{"-i", hostFile, playbookFile.Name()}
	result, err = utils.ExecShell("ansible-playbook", commandArgs...)

	return
}

// 判断 获取临时目录，默认/dev/shm
func GetTmpDir() string {
	dir := "/dev/shm"
	_, err := os.Stat(dir)
	if os.IsPermission(err) {
		return ""
	}

	return dir
}
