package batch

import (
	"time"

	"sync"

	"log"
)

var instance *TimeoutBatcher
var once sync.Once

type TaskFunc func() error

func GetTimeOutBatcherInstance() *TimeoutBatcher {

	once.Do(func() {
		instance = newTimeoutBatcher()
	})

	return instance
}

type Task struct {
	T        *time.Ticker
	quit     chan struct{}
	taskFunc func() error
}

func NewTimer(duration time.Duration, taskFunc func() error) Task {
	return Task{
		quit:     make(chan struct{}, 1),
		T:        time.NewTicker(duration),
		taskFunc: taskFunc,
	}
}

func (t *Task) Start() error {

	for {
		select {
		case <-t.T.C:
			if err := t.taskFunc(); err != nil {
				t.quit <- struct{}{}
			}
		case <-t.quit:
			t.Stop()
			return nil
		}
	}

	return nil
}

func (t *Task) Stop() {
	t.T.Stop()
}

type TimeoutBatcher struct {
	timers map[string]Task
}

func newTimeoutBatcher() *TimeoutBatcher {

	return &TimeoutBatcher{
		timers: make(map[string]Task),
	}
}

func (t *TimeoutBatcher) Run(taskFunc TaskFunc, duration time.Duration) chan struct{} {

	timer := NewTimer(duration, taskFunc)

	var err error

	go func() {
		defer log.Println("timer is closing")
		err = timer.Start()

		if err != nil {
			log.Println(err)
			return
		}
	}()

	return timer.quit
}
