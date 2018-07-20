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
	"sync"
	"testing"
	"time"

	"github.com/it-chain/engine/common/amqp/pubsub"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

var wg sync.WaitGroup

func TestConnect(t *testing.T) {

	wg.Add(2)
	c := pubsub.Connect("")
	defer c.Close()

	handler := &MockHandler{}
	handler.HandleNameUpdateCommandFunc = func(event UserNameUpdateEvent) {
		assert.Equal(t, event.Name, "JUN")
		assert.Equal(t, event.ID, "123")
		wg.Done()
	}

	handler.HandleFunc = func(command UserAddCommand) {
		assert.Equal(t, command.ID, "123")
		wg.Done()
	}

	err := c.Subscribe("asd", "asd", handler)
	assert.NoError(t, err)

	err = c.Publish("asd", "asd", UserNameUpdateEvent{
		Name: "JUN",
		EventModel: midgard.EventModel{
			ID:      "123",
			Time:    time.Now(),
			Type:    "123",
			Version: 1,
		}})
	assert.NoError(t, err)

	err = c.Publish("asd", "asd", UserAddCommand{
		CommandModel: midgard.CommandModel{
			ID: "123",
		}})

	assert.NoError(t, err)

	wg.Wait()
}

type UserNameUpdateEvent struct {
	midgard.EventModel
	Name string
}

type UserAddCommand struct {
	midgard.CommandModel
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
