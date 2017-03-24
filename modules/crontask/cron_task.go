package crontask

import (
	"github.com/robfig/cron"
	"errors"
	"sync"
	"strings"
)

var DefaultCronTask *CronTask

type CronMap map[string]*cron.Cron

type CronTask struct {
	sync.RWMutex
	tasks CronMap
}

func NewCronTask() *CronTask {
	return &CronTask {
		sync.RWMutex{},
		make(CronMap),
	}
}

// 新增定时任务,如果name存在，则添加失败
// name 任务名称
// spec crontab时间格式定义  可定义多个时间\n分隔
func(cronTask *CronTask) Add(name string, spec string, cmd cron.FuncJob ) (err error) {
	if name == "" || spec == "" || cmd == nil {
		return errors.New("参数不完整")
	}
	if cronTask.IsExist(name) {
		return errors.New("任务已存在")
	}

	spec = strings.TrimSpace(spec)
	cronTask.Lock()
	defer cronTask.Unlock()
	cronTask.tasks[name] = cron.New()
	specs := strings.Split(spec, "|||")
	for _, item := range(specs) {
		_, err = cron.Parse(item)
		if err != nil {
			return err
		}
	}
	for _, item := range(specs) {
		err = cronTask.tasks[name].AddFunc(item, cmd)
	}

	return err
}

// 任务不存在则新增，任务已存在则删除后新增
func(cronTask *CronTask) AddOrReplace(name string, spec string, cmd cron.FuncJob) error {
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
func(cronTask *CronTask) Run() {
	for _, cron := range cronTask.tasks {
		// cron内部有开启goroutine,此处不用新建goroutine
		cron.Start()
	}
}