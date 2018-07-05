package gateway_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

type MockRepostiory struct {
	loadFunc func(aggregate midgard.Aggregate, aggregateID string) error
	saveFunc func(aggregateID string, events ...midgard.Event) error
}

func (m MockRepostiory) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return m.loadFunc(aggregate, aggregateID)
}

func (m MockRepostiory) Save(aggregateID string, events ...midgard.Event) error {
	return m.saveFunc(aggregateID, events...)
}

func (MockRepostiory) Close() {

}

func TestNewConnection(t *testing.T) {

	c := gateway.Connection{
		Address: "127.0.0.1",
		AggregateModel: midgard.AggregateModel{
			ID: "conn1",
		},
	}

	repo := MockRepostiory{}
	repo.saveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, aggregateID, "conn1")
		assert.Equal(t, len(events), 1)
		assert.IsType(t, &gateway.ConnectionCreatedEvent{}, events[0])
		return nil
	}

	eventstore.InitForMock(repo)

	conn, err := gateway.NewConnection(c.ID, c.Address)
	assert.NoError(t, err)

	assert.Equal(t, conn, c)
}
