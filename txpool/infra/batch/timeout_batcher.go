package batch

import (
	"time"

	"sync"

	"log"
)

var instance *TimeoutBatcher
var once sync.Once

type TimerFunc func() error

func GetTimeOutBatcherInstance() *TimeoutBatcher {

	once.Do(func() {
		instance = newTimeoutBatcher()
	})

	return instance
}

type Timer struct {
	T         *time.Ticker
	quit      chan struct{}
	timerFunc func() error
}

func NewTimer(duration time.Duration, timerFunc func() error) Timer {
	return Timer{
		quit:      make(chan struct{}, 1),
		T:         time.NewTicker(duration),
		timerFunc: timerFunc,
	}
}

func (t *Timer) Start() error {

	for {
		select {
		case <-t.T.C:
			if err := t.timerFunc(); err != nil {
				t.quit <- struct{}{}
			}
		case <-t.quit:
			t.Stop()
			return nil
		}
	}

	return nil
}

func (t *Timer) Stop() {
	t.T.Stop()
}

type TimeoutBatcher struct {
	timers map[string]Timer
}

func newTimeoutBatcher() *TimeoutBatcher {

	return &TimeoutBatcher{
		timers: make(map[string]Timer),
	}
}

func (t *TimeoutBatcher) Register(timerFunc TimerFunc, duration time.Duration) chan struct{} {

	timer := NewTimer(duration, timerFunc)

	var err error

	go func() {
		defer log.Println("timer is closing")
		err = timer.Start()

		if err != nil {
			return
		}
	}()

	return timer.quit
}
