package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/api"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

type MockEventRepository struct {
	SaveFunc func(aggregateID string, events ...midgard.Event) error
}

func (rp MockEventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return nil
}

func (rp MockEventRepository) Save(aggregateID string, events ...midgard.Event) error {
	return rp.SaveFunc(aggregateID, events...)
}

func (rp MockEventRepository) Close() {}

func TestTransactionApi_CreateTransaction(t *testing.T) {
	tests := map[string]struct {
		input struct {
			txID   string
			txData txpool.TxData
		}
		err error
	}{
		"success": {
			input: struct {
				txID   string
				txData txpool.TxData
			}{txID: "zf", txData: txpool.TxData{ID: "gg"}},
			err: nil,
		},
	}

	eventRepository := MockEventRepository{}
	eventRepository.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, "zf", aggregateID)
		assert.Equal(t, "zf", events[0].GetID())
		return nil
	}

	transactionApi := api.NewTransactionApi(eventRepository, "zf")

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := transactionApi.CreateTransaction(test.input.txID, test.input.txData)

		assert.Equal(t, test.err, err)
	}
}
