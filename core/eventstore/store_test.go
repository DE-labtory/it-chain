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

package eventstore_test

import (
	"os"
	"testing"

	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestSaveAndLoad(t *testing.T) {

	defer InitStore()()

	aggregateID := "123"
	event := UserCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   aggregateID,
			Type: "User",
		},
	}

	err := eventstore.Save(event.GetID(), event)
	assert.NoError(t, err)

	user := &User{}

	err = eventstore.Load(user, "123")
	assert.NoError(t, err)

	assert.Equal(t, user.ID, aggregateID)
}

func TestPanicWhenInitTwice(t *testing.T) {

	assert.Panics(t, func() {
		defer InitStore()()
		defer InitStore()()
	})
}

func InitStore() func() {

	path := "./.test"
	eventstore.InitLevelDBStore(path, nil,
		UserCreatedEvent{},
		UserNameUpdatedEvent{})

	return func() {
		eventstore.Close()
		os.RemoveAll(path)
	}
}
