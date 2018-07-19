package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/magiconair/properties/assert"
)

func TestBlockApi_AddBlockToPool(t *testing.T) {
	tests := map[string]struct {
		input struct {
			block blockchain.Block
		}
	}{
		"success": {
			input: struct {
				block blockchain.Block
			}{block: &blockchain.DefaultBlock{
				Height: uint64(11),
			}},
		},
	}

	publisherId := "zf"

	blockApi, _ := api.NewBlockApi(publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi.AddBlockToPool(test.input.block)
	}
}

func TestBlockApi_CheckAndSaveBlockFromPool(t *testing.T) {
	tests := map[string]struct {
		input struct {
			height blockchain.BlockHeight
		}
		err error
	}{
		"success": {
			input: struct {
				height blockchain.BlockHeight
			}{height: blockchain.BlockHeight(12)},
			err: nil,
		},
	}
	publisherId := "zf"

	// When
	blockApi, _ := api.NewBlockApi(publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		err := blockApi.CheckAndSaveBlockFromPool(test.input.height)

		// Then
		assert.Equal(t, test.err, err)
	}
}

func TestBlockApi_SyncIsProgressing(t *testing.T) {
	// when
	publisherId := "zf"

	// when
	blockApi, _ := api.NewBlockApi(publisherId)

	// then
	state := blockApi.SyncIsProgressing()
	assert.Equal(t, blockchain.DONE, state)
}
