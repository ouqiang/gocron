package timewheel

import (
    "time"
    "container/ring"
    "container/list"
)

// @author qiang.ou<qingqianludao@gmail.com>

type TimeWheel struct {
    interval time.Duration
    ticker *time.Ticker
    slots *ring.Ring
    slotNum int
    taskChannel chan Task
    stopChannel chan bool
}


type Task struct {
    delay time.Duration
    circle int
    job Job
}

type Job func()

func New(interval time.Duration, slotNum int) *TimeWheel {
    if  interval <= 0 || slotNum <= 0 {
        return nil
    }
    tw := &TimeWheel{
        interval: interval,
        slots: ring.New(slotNum),
        slotNum: slotNum,
        taskChannel: make(chan Task),
        stopChannel: make(chan bool),
    }

    tw.initSlots()

    return tw
}

func (tw *TimeWheel) initSlots()  {
    for i := 0; i < tw.slots.Len(); i++ {
        tw.slots.Value = list.New()
        tw.slots = tw.slots.Next()
    }
}

func (tw *TimeWheel) Start()  {
    tw.ticker = time.NewTicker(tw.interval)
    go tw.start()
}

func (tw *TimeWheel) Add(delay time.Duration, job Job)  {
    if delay < 0 || job == nil {
        return
    }
    tw.taskChannel <- Task{delay:delay, job: job}
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
    l := tw.slots.Value.(*list.List)
    tw.scanAndRunTask(l)
    tw.slots = tw.slots.Next()
}

func (tw *TimeWheel) scanAndRunTask(l *list.List)  {
    for e := l.Front(); e != nil; {
        task := e.Value.(*Task)
        if task.circle > 0 {
            task.circle--
            e = e.Next()
            continue
        }

        go task.job()
        next := e.Next()
        l.Remove(e)
        e = next
    }
}

func (tw *TimeWheel) addTask(task *Task)  {
    step, circle := tw.getStepAndCircle(task.delay)
    task.circle = circle

    l := tw.slots.Move(step).Value.(*list.List)
    l.PushBack(task)
}

func (tw *TimeWheel) getStepAndCircle(d time.Duration) (step int, circle int) {
    delaySeconds := int(d.Seconds())
    intervalSeconds := int(tw.interval.Seconds())
    circle = int(delaySeconds / intervalSeconds /  tw.slotNum)
    step = int(delaySeconds / intervalSeconds) % tw.slotNum

    return
}