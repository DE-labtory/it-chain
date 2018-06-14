package batch

import (
	"time"

	"sync"

	"github.com/rs/xid"
)

var instance *TimeoutBatcher
var once sync.Once

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
		quit:      make(chan struct{}),
		T:         time.NewTicker(duration),
		timerFunc: timerFunc,
	}
}

func (t *Timer) Start() error {

	for {
		select {
		case <-t.T.C:
			if err := t.timerFunc(); err != nil {
				return err
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

func (t *TimeoutBatcher) Register(timerFunc func() error, duration time.Duration) chan struct{} {

	timer := NewTimer(duration, timerFunc)

	var err error

	go func() {
		err = timer.Start()

		if err != nil {
			return
		}
	}()

	t.timers[xid.New().String()] = timer

	return timer.quit
}

func (t *TimeoutBatcher) StopAll() {
	for key, timer := range t.timers {
		timer.Stop()
		delete(t.timers, key)
	}
}
