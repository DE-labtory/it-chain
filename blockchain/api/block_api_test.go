package api_test

import (
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"testing"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/magiconair/properties/assert"
)

type MockBlockRepository struct {}

func (br MockBlockRepository) Close() {}
func (br MockBlockRepository) GetValidator() common.Validator { return nil }
func (br MockBlockRepository) AddBlock(block common.Block) error { return nil }
func (br MockBlockRepository) GetLastBlock(block common.Block) error { return nil }
func (br MockBlockRepository) NewEmptyBlock() (blockchain.Block, error) {return nil, nil}


type MockBlockPool struct {
	AddFunc func(block blockchain.Block)
	GetFunc func(height blockchain.BlockHeight) blockchain.Block
	blockPool map[blockchain.BlockHeight] blockchain.Block
}

func (bp MockBlockPool) Add(block blockchain.Block) {
	bp.AddFunc(block)
}
func (bp MockBlockPool) Get(height blockchain.BlockHeight) blockchain.Block {
	return bp.GetFunc(height)
}
func (bp MockBlockPool) Delete(height blockchain.Block) {}

func TestBlockApi_AddBlockToPool(t *testing.T) {
	tests := map[string] struct {
		input struct {
			block blockchain.Block
		}
	} {
		"success": {
			input: struct {
				block blockchain.Block
			} {block: &blockchain.DefaultBlock{
				Height: uint64(11),
			}},
		},
		"block nil test": {
			input: struct {
				block blockchain.Block
			} {block: nil},
		},
	}

	blockRepository := MockBlockRepository{}
	publisherId := "zf"
	blockPool := MockBlockPool{}
	blockPool.AddFunc = func(block blockchain.Block) {
		assert.Equal(t, uint64(11), block.GetHeight())
	}

	blockApi, _ := api.NewBlockApi(blockRepository, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi.AddBlockToPool(test.input.block)
	}
}


func TestBlockApi_CheckAndSaveBlockFromPool(t *testing.T) {
	tests := map[string] struct {
		input struct {
			block blockchain.Block
		}
		err error
	} {
		"success": {
			input: struct {
				block blockchain.Block
			}{block: &blockchain.DefaultBlock{
				Height: blockchain.BlockHeight(12),
			}},
			err: nil,
		},
		"block nil test": {
			input: struct {
				block blockchain.Block
			}{block: &blockchain.DefaultBlock{
				Height: blockchain.BlockHeight(13),
			}},
			err: api.ErrNilBlock,
		},
	}
	// When
	blockRepository := MockBlockRepository{}
	publisherId := "zf"

	// predefined block pool
	blockPool := MockBlockPool{
		blockPool: map[blockchain.BlockHeight] blockchain.Block {
			12: &blockchain.DefaultBlock{
				Height: blockchain.BlockHeight(12),
			},
		},
	}
	blockPool.GetFunc = func(height blockchain.BlockHeight) blockchain.Block {
		return blockPool.blockPool[height]
	}
	// When
	blockApi, _ := api.NewBlockApi(blockRepository, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		err := blockApi.CheckAndSaveBlockFromPool(test.input.block)

		// Then
		assert.Equal(t, test.err, err)
	}
}