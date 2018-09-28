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

package pubsub_test

import (
	"testing"

	"encoding/json"
	"reflect"

	"sync"

	"fmt"

	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/midgard"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestTopicSubscriber_SubscribeTopic(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(8)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	handler1 := &MockHandler{}
	handler1.HandleNameUpdateCommandFunc = func(event UserNameUpdateEvent) {
		fmt.Println("1")
		assert.Equal(t, event.Name, "Jun")
		wg.Done()
	}

	handler1.HandleFunc = func(command UserAddCommand) {
		fmt.Println("2")
		assert.Equal(t, command.Age, 123)
		wg.Done()
	}

	handler2 := &MockHandler{}
	handler2.HandleNameUpdateCommandFunc = func(event UserNameUpdateEvent) {
		fmt.Println("1")
		assert.Equal(t, event.Name, "Jun")
		wg.Done()
	}

	handler2.HandleFunc = func(command UserAddCommand) {
		fmt.Println("2")
		assert.Equal(t, command.Age, 123)
		wg.Done()
	}

	handler3 := &MockHandler{}
	handler3.HandleNameUpdateCommandFunc = func(event UserNameUpdateEvent) {
		assert.Equal(t, event.Name, "Jun")
		wg.Done()
	}

	handler3.HandleFunc = func(command UserAddCommand) {
		assert.Equal(t, command.Age, 123)
		wg.Done()
	}

	subscriber.SubscribeTopic("test.*", handler1)
	subscriber.SubscribeTopic("test.*", handler2)
	//subscriber.SubscribeTopic("test.*", handler3)

	publish(t, "Event", "test.created", UserNameUpdateEvent{Name: "Jun"})
	publish(t, "Event", "test.created", UserAddCommand{Age: 123})
	publish(t, "Event", "test.created", UserNameUpdateEvent{Name: "Jun"})
	publish(t, "Event", "test.created", UserNameUpdateEvent{Name: "Jun"})
	publish(t, "Event", "test.created", UserNameUpdateEvent{Name: "Jun"})
	publish(t, "Event", "test.created", UserNameUpdateEvent{Name: "Jun"})
	publish(t, "Event", "test.created", UserNameUpdateEvent{Name: "Jun"})
	publish(t, "Event", "test.created", UserNameUpdateEvent{Name: "Jun"})
	publish(t, "Event", "test.created", UserNameUpdateEvent{Name: "Jun"})

	wg.Wait()
}

type UserNameUpdateEvent struct {
	midgard.EventModel
	Name string
}

type UserAddCommand struct {
	midgard.CommandModel
	Age int
}

type MockHandler struct {
	HandleFunc                  func(command UserAddCommand)
	HandleNameUpdateCommandFunc func(event UserNameUpdateEvent)
}

func (d *MockHandler) Handle(command UserAddCommand) {
	d.HandleFunc(command)
}

func (d *MockHandler) HandleNameUpdateCommand(event UserNameUpdateEvent) {
	d.HandleNameUpdateCommandFunc(event)
}

func publish(t *testing.T, exchange string, topic string, data interface{}) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	assert.NoError(t, err)

	ch, err := conn.Channel()
	assert.NoError(t, err)

	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	assert.NoError(t, err)

	var matchingValue string

	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		matchingValue = reflect.TypeOf(data).Elem().Name()
	} else {
		matchingValue = reflect.TypeOf(data).Name()
	}

	b, err := json.Marshal(data)

	message := pubsub.Message{
		MatchingValue: matchingValue,
		Data:          b,
	}

	byte, err := json.Marshal(message)
	assert.NoError(t, err)

	err = ch.Publish(
		exchange, // exchange
		topic,    // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        byte,
		})

	assert.NoError(t, err)
}
