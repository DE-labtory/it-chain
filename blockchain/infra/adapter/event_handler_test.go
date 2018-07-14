package adapter_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"time"
	"github.com/it-chain/it-chain-Engine/blockchain/test/mock"
	"github.com/magiconair/properties/assert"
)

func TestEventHandler_HandleBlockAddToPoolEvent(t *testing.T) {
	tests := map[string]struct {
		input struct {
			blockchain.BlockAddToPoolEvent
		}
		err error
	}{
		"success": {
			input: struct {
				blockchain.BlockAddToPoolEvent
			}{BlockAddToPoolEvent: blockchain.BlockAddToPoolEvent{
				Seal: []byte{0x1},
				PrevSeal: []byte{0x1},
				Height: uint64(12),
				TxList: []byte{0x1},
				TxSeal: [][]byte{{0x1}},
				Timestamp: time.Now(),
				Creator: []byte{0x1},
			}},
			err: nil,
		},
		"event w/o height test": {
			input: struct {
				blockchain.BlockAddToPoolEvent
			}{BlockAddToPoolEvent: blockchain.BlockAddToPoolEvent{
				Seal: []byte{0x1},
				PrevSeal: []byte{0x1},
				TxList: []byte{0x1},
				TxSeal: [][]byte{{0x1}},
				Timestamp: time.Now(),
				Creator: []byte{0x1},
			}},
			err: adapter.ErrBlockMissingProperties,
		},
		"event w/o seal, prevseal test": {
			input: struct {
				blockchain.BlockAddToPoolEvent
			}{BlockAddToPoolEvent: blockchain.BlockAddToPoolEvent{
				TxList: []byte{0x1},
				TxSeal: [][]byte{{0x1}},
				Timestamp: time.Now(),
				Creator: []byte{0x1},
			}},
			err: adapter.ErrBlockMissingProperties,
		},

	}

	blockApi := mock.MockBlockApi{}
	blockApi.CheckAndSaveBlockFromPoolFunc = func(height blockchain.BlockHeight) error {
		assert.Equal(t, height, uint64(12))
		return nil
	}
	eventHandler := adapter.NewEventHandler(blockApi)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := eventHandler.HandleBlockAddToPoolEvent(test.input.BlockAddToPoolEvent)

		assert.Equal(t, err, test.err)
	}
}
