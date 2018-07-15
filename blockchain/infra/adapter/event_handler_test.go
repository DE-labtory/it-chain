package adapter_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/it-chain/it-chain-Engine/blockchain/test/mock"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

func TestEventHandler_HandleBlockAddToPoolEvent(t *testing.T) {
	tests := map[string]struct {
		input struct {
			blockchain.BlockStagedEvent
		}
		err error
	}{
		"success": {
			input: struct {
				blockchain.BlockStagedEvent
			}{BlockStagedEvent: blockchain.BlockStagedEvent{
				EventModel: midgard.EventModel{
					ID: "zf",
				},
				State: blockchain.Staged,
			}},
			err: nil,
		},
	}

	blockApi := mock.BlockApi{}
	blockApi.CommitBlockFromPoolOrSyncFunc = func(blockId string) error {
		assert.Equal(t, blockId, "zf")
		return nil
	}
	eventHandler := adapter.NewEventHandler(blockApi)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := eventHandler.HandleBlockAddToPoolEvent(test.input.BlockStagedEvent)

		assert.Equal(t, err, test.err)
	}
}
