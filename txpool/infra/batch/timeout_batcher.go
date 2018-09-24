/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package batch

import (
	"time"

	"sync"

	"github.com/it-chain/iLogger"
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
				iLogger.Error(nil, "error: "+err.Error())
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
		//defer log.Println("timer is closing")
		err = timer.Start()

		if err != nil {
			iLogger.Error(nil, "error: "+err.Error())
			//	return
		}
	}()

	return timer.quit
}
