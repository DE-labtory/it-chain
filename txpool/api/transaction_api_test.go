package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/api"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

type MockEventRepository struct {
	SaveFunc func(aggregateID string, events ...midgard.Event) error
	LoadFunc func(aggregate midgard.Aggregate, aggregateID string) error
}

func (rp MockEventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return rp.LoadFunc(aggregate, aggregateID)
}

func (rp MockEventRepository) Save(aggregateID string, events ...midgard.Event) error {
	return rp.SaveFunc(aggregateID, events...)
}

func (rp MockEventRepository) Close() {}

func TestTransactionApi_CreateTransaction(t *testing.T) {

	tests := map[string]struct {
		input struct {
			txData txpool.TxData
		}
		err error
	}{
		"success": {
			input: struct {
				txData txpool.TxData
			}{txData: txpool.TxData{ICodeID: "gg"}},
			err: nil,
		},
	}

	eventRepository := MockEventRepository{}
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, "gg", events[0].(*txpool.TxCreatedEvent).ICodeID)
		return nil
	}

	eventstore.InitForMock(eventRepository)

	transactionApi := api.NewTransactionApi("zf")

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		_, err := transactionApi.CreateTransaction(test.input.txData)

		assert.Equal(t, test.err, err)
	}
}

func TestTransactionApi_DeleteTransaction(t *testing.T) {

	tests := map[string]struct {
		input string
		err   error
	}{
		"success": {
			input: "transactionID",
			err:   nil,
		},
	}

	eventRepository := MockEventRepository{}
	eventRepository.LoadFunc = func(aggregate midgard.Aggregate, aggregateID string) error {

		aggregate.(*txpool.Transaction).TxId = "transactionID"
		return nil
	}

	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {

		assert.Equal(t, "transactionID", events[0].(*txpool.TxDeletedEvent).GetID())
		return nil
	}

	eventstore.InitForMock(eventRepository)

	transactionApi := api.NewTransactionApi("zf")

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := transactionApi.DeleteTransaction(test.input)

		assert.Equal(t, test.err, err)
	}
}
