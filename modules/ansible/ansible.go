package ansible

// ansible ad-hoc 命令封装

import (
	"errors"
	"github.com/ouqiang/cron-scheduler/modules/utils"
)



/**
 * 执行ad-hoc
 * hosts  主机名或主机别名 逗号分隔
 * hostFile 主机名文件
 * module 调用模块
 * args   传递给模块的参数
*/
func ExecCommand(hosts string, hostFile string, args... string) (output string, err error) {
	if hosts== "" || hostFile == "" || len(args) == 0 {
		err = errors.New("参数不完整")
		return
	}
	commandArgs := []string{hosts, "-i",  hostFile}
	commandArgs = append(commandArgs,  args...)
	output, err = utils.ExecShell("ansible", commandArgs...)

	return
}