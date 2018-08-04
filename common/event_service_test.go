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

package common_test

import (
	"sync"
	"testing"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/stretchr/testify/assert"
)

type SampleEvent struct {
	Attr string
}

type MockHandler struct {
	HandleFunc func(event SampleEvent)
}

func (d *MockHandler) Handle(event SampleEvent) {
	d.HandleFunc(event)
}

func TestEventServiceImpl_Publish(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	handler := &MockHandler{}
	handler.HandleFunc = func(event SampleEvent) {
		assert.Equal(t, "test", event.Attr)
		wg.Done()
	}

	subscriber.SubscribeTopic("test.*", handler)

	eventService := common.NewEventService("", "Event")

	err := eventService.Publish("test.sample", SampleEvent{Attr: "test"})
	assert.NoError(t, err)
	wg.Wait()
}
