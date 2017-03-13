package utils

import (
	"github.com/robfig/cron"
	"errors"
	"scheduler/utils/app"
	"sync"
)

var DefaultCronTask CronTask;

type CronTask struct {
	sync.RWMutex
	tasks map[string]*cron.Cron
}

func init()  {
	if app.Installed {
		DefaultCronTask = CronTask{
			sync.RWMutex{},
			make(map[string]*cron.Cron),
		}
	}
}

// 新增定时任务,如果name存在，则添加失败
func(cronTask *CronTask) Add(name string, spec string, cmd func() ) error {
	if name == "" || spec == "" || cmd == nil {
		return errors.New("参数不完整")
	}
	if cronTask.IsExist(name) {
		return errors.New("任务已存在")
	}

	cronTask.Lock()
	defer cronTask.Unlock()
	cronTask.tasks[name] = cron.New()
	err := cronTask.tasks[name].AddFunc(spec, cmd)

	return err
}

// 任务不存在则新增，任务已存在则替换任务
func(cronTask *CronTask) addOrReplace(name string, spec string, cmd func() ) error {
	if cronTask.IsExist(name) {
		cronTask.Delete(name)
	}

	return cronTask.Add(name, spec, cmd)
}


// 判断任务是否存在
func(cronTask *CronTask) IsExist(name string) bool {
	cronTask.RLock()
	defer cronTask.RUnlock()
	_, ok := cronTask.tasks[name]

	return ok
}

// 启动任务
func(cronTask *CronTask) Start(name string) {
	if cronTask.IsExist(name) {
		cronTask.tasks[name].Start()
	}
}

// 停止任务
func(cronTask *CronTask) Stop(name string) {
	if cronTask.IsExist(name) {
		cronTask.tasks[name].Stop()
	}
}

// 删除任务
func(cronTask *CronTask) Delete(name string) {
	cronTask.Stop(name)
	cronTask.Lock()
	defer cronTask.Unlock()
	delete(cronTask.tasks, name)
}

// 运行所有任务
func(cronTask *CronTask) run() {
	for _, cron := range cronTask.tasks {
		// cron内部有开启goroutine,此处不用新建goroutine
		cron.Start()
	}
}