/*
 * Copyright 2018 DE-labtory
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

package common

import (
	"bytes"
	"errors"
	"reflect"
	"runtime"
	"strconv"

	"github.com/DE-labtory/engine/common/rabbitmq/pubsub"
)

var ErrEventType = errors.New("Error type of event is not struct")

type EventService interface {
	Publish(topic string, event interface{}) error
	Close()
}

type EventServiceImpl struct {
	mqUrl      string
	exchange   string
	channelMap map[uint64]worker
}

func NewEventService(mqUrl string, exchange string) *EventServiceImpl {
	return &EventServiceImpl{
		mqUrl:      mqUrl,
		exchange:   exchange,
		channelMap: make(map[uint64]worker),
	}
}

func (s *EventServiceImpl) Publish(topic string, event interface{}) error {
	if !eventIsStruct(event) {
		return ErrEventType
	}

	w, ok := s.channelMap[getGID()]
	if !ok {
		worker := NewWorker(make(chan message, 1), pubsub.NewTopicPublisher(s.mqUrl, s.exchange))
		s.channelMap[getGID()] = worker
		go worker.work()
		worker.message <- message{
			event: event,
			topic: topic,
		}

		return nil
	}

	w.message <- message{
		event: event,
		topic: topic,
	}

	return nil
}

func (s *EventServiceImpl) Close() {
	for _, w := range s.channelMap {
		w.quit <- struct{}{}
	}
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func eventIsStruct(event interface{}) bool {
	return reflect.TypeOf(event).Kind() == reflect.Struct
}

type worker struct {
	quit           chan struct{}
	message        chan message
	topicPublisher *pubsub.TopicPublisher
}

type message struct {
	topic string
	event interface{}
}

func NewWorker(chMessage chan message, topicPublisher *pubsub.TopicPublisher) worker {
	return worker{
		quit:           make(chan struct{}),
		message:        chMessage,
		topicPublisher: topicPublisher,
	}
}

func (w *worker) work() {
	for {
		select {
		case m := <-w.message:
			w.topicPublisher.Publish(m.topic, m.event)
		case <-w.quit:
			return
		}
	}
}
