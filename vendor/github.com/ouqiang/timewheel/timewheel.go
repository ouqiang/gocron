package timewheel

import (
    "time"
    "container/list"
)

// @author qiang.ou<qingqianludao@gmail.com>

type Job func([]interface{})

type TimeWheel struct {
    interval time.Duration
    ticker *time.Ticker
    slots []*list.List
    currentPos int
    slotNum int
    job Job
    taskChannel chan Task
    stopChannel chan bool
}


type Task struct {
    delay time.Duration
    circle int
    data []interface{}
}

func New(interval time.Duration, slotNum int, job Job) *TimeWheel {
    if  interval <= 0 || slotNum <= 0 || job == nil {
        return nil
    }
    tw := &TimeWheel{
        interval: interval,
        slots: make([]*list.List, slotNum),
        currentPos: 0,
        job: job,
        slotNum: slotNum,
        taskChannel: make(chan Task),
        stopChannel: make(chan bool),
    }

    tw.initSlots()

    return tw
}

func (tw *TimeWheel) initSlots()  {
    for i := 0; i < tw.slotNum; i++ {
        tw.slots[i] = list.New()
    }
}

func (tw *TimeWheel) Start()  {
    tw.ticker = time.NewTicker(tw.interval)
    go tw.start()
}

func (tw *TimeWheel) Add(delay time.Duration, data []interface{})  {
    if delay <= 0  {
        return
    }
    tw.taskChannel <- Task{delay:delay, data: data}
}

func (tw *TimeWheel) Stop()  {
    tw.stopChannel <- true
}

func (tw *TimeWheel) start()  {
    for {
        select {
        case <- tw.ticker.C:
            tw.tickHandler()
        case task := <- tw.taskChannel:
            tw.addTask(&task)
        case <- tw.stopChannel:
            tw.ticker.Stop()
            return
        }
    }
}

func (tw *TimeWheel) tickHandler()  {
    l := tw.slots[tw.currentPos]
    tw.scanAndRunTask(l)
    if tw.currentPos == tw.slotNum - 1 {
        tw.currentPos = 0
    } else {
        tw.currentPos++
    }
}

func (tw *TimeWheel) scanAndRunTask(l *list.List)  {
    for e := l.Front(); e != nil; {
        task := e.Value.(*Task)
        if task.circle > 0 {
            task.circle--
            e = e.Next()
            continue
        }

        go tw.job(task.data)
        next := e.Next()
        l.Remove(e)
        e = next
    }
}

func (tw *TimeWheel) addTask(task *Task)  {
    pos, circle := tw.getPositionAndCircle(task.delay)
    task.circle = circle

    tw.slots[pos].PushBack(task)
}

func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
    delaySeconds := int(d.Seconds())
    intervalSeconds := int(tw.interval.Seconds())
    circle = int(delaySeconds / intervalSeconds /  tw.slotNum)
    pos = int(tw.currentPos + delaySeconds / intervalSeconds) % tw.slotNum


    return
}