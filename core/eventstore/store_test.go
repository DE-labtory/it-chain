package eventstore_test

import (
	"testing"

	"os"

	"fmt"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestEventStore(t *testing.T) {

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

	fmt.Println(eventstore.Instance)

	user := &User{}

	err = eventstore.Load(user, "123")
	assert.NoError(t, err)

	assert.Equal(t, user.ID, aggregateID)
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
