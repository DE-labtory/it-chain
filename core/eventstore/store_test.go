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
