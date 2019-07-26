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

package pubsub_test

import (
	"sync"
	"testing"

	"github.com/DE-labtory/it-chain/common/rabbitmq/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestTopicPublisher_Publish(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(2)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	publisher := pubsub.NewTopicPublisher("", "Event")

	defer subscriber.Close()

	handler := &MockHandler{}
	handler.HandleNameUpdateCommandFunc = func(event UserNameUpdateEvent) {
		assert.Equal(t, event.Name, "Jun")
		wg.Done()
	}

	handler.HandleFunc = func(command UserAddCommand) {
		assert.Equal(t, command.Age, 123)
		wg.Done()
	}

	subscriber.SubscribeTopic("test.*", handler)

	publisher.Publish("test.created", UserNameUpdateEvent{Name: "Jun"})
	publisher.Publish("test.created", UserAddCommand{Age: 123})
	wg.Wait()
}
