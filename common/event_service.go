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

package common

import (
	"errors"
	"reflect"

	"github.com/it-chain/engine/common/rabbitmq/pubsub"
)

var ErrEventType = errors.New("Error type of event is not struct")

type EventService interface {
	Publish(topic string, event interface{}) error
}

type EventServiceImpl struct {
	topicPublisher pubsub.TopicPublisher
}

func NewEventService(mqUrl string, exchange string) EventService {
	return &EventServiceImpl{
		topicPublisher: pubsub.NewTopicPublisher(mqUrl, exchange),
	}
}

func (s *EventServiceImpl) Publish(topic string, event interface{}) error {
	if !eventIsStruct(event) {
		return ErrEventType
	}

	err := s.topicPublisher.Publish(topic, event)
	if err != nil {
		return err
	}

	return nil
}

func eventIsStruct(event interface{}) bool {
	return reflect.TypeOf(event).Kind() == reflect.Struct
}
