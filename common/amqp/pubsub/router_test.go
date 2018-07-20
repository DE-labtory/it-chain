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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/it-chain/engine/common/amqp/pubsub"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestNewParamBasedRouter(t *testing.T) {
	d, err := pubsub.NewParamBasedRouter()
	assert.NoError(t, err)

	handler := &Handler{}
	handler.HandleNameUpdateCommandFunc = func(command UserNameUpdateCommand) {
		assert.Equal(t, command.Name, "jun")
	}

	err = d.SetHandler(handler)
	assert.NoError(t, err)

	cmd := UserNameUpdateCommand{
		Name: "jun",
	}

	b, err := json.Marshal(cmd)
	assert.NoError(t, err)

	fmt.Println(b)

	err = d.Route(b, "UserNameUpdateCommand")
	assert.NoError(t, err)
}

type UserNameUpdateCommand struct {
	midgard.EventModel
	Name string
}

type Handler struct {
	HandleNameUpdateCommandFunc func(command UserNameUpdateCommand)
}

func (d *Handler) HandleNameUpdateCommand(command UserNameUpdateCommand) {
	d.HandleNameUpdateCommandFunc(command)
}
