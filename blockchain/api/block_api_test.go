package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/it-chain/it-chain-Engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestBlockApi_StageBlock(t *testing.T) {
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

	blockQueryApi := mock.BlockQueryApi{}
	publisherId := "zf"

	blockApi, _ := api.NewBlockApi(blockQueryApi, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi.StageBlock(test.input.block)
	}
}

func TestBlockApi_CommitBlockFromPoolOrSync(t *testing.T) {
	tests := map[string]struct {
		input struct {
			height  blockchain.BlockHeight
			blockId string
		}
		err error
	}{
		"success": {
			input: struct {
				height  blockchain.BlockHeight
				blockId string
			}{height: blockchain.BlockHeight(12), blockId: "zf"},
			err: nil,
		},
	}
	// When
	blockQueryApi := mock.BlockQueryApi{}
	blockQueryApi.GetStagedBlockByIdFunc = func(blockId string) (blockchain.Block, error) {
		assert.Equal(t, "zf", blockId)

		return &blockchain.DefaultBlock{
			Height: blockchain.BlockHeight(12),
		}, nil
	}
	blockQueryApi.GetLastCommitedBlockFunc = func() (blockchain.Block, error) {
		return &blockchain.DefaultBlock{
			Height: blockchain.BlockHeight(12),
		}, nil
	}
	publisherId := "zf"

	// When
	blockApi, _ := api.NewBlockApi(blockQueryApi, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		err := blockApi.CommitBlockFromPoolOrSync(test.input.blockId)

		// Then
		assert.Equal(t, test.err, err)
	}
}

func TestBlockApi_SyncedCheck(t *testing.T) {
	// TODO:
}

func TestBlockApi_SyncIsProgressing(t *testing.T) {
	// when
	blockQueryApi := mock.BlockQueryApi{}
	publisherId := "zf"

	// when
	blockApi, _ := api.NewBlockApi(blockQueryApi, publisherId)

	// then
	state := blockApi.SyncIsProgressing()
	assert.Equal(t, blockchain.DONE, state)
}
