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
}

func (bp MockBlockPool) Add(block blockchain.Block) {
	bp.AddFunc(block)
}
func (bp MockBlockPool) Get(height blockchain.BlockHeight) blockchain.Block { return nil }
func (bp MockBlockPool) Delete(height blockchain.BlockHeight) {}

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

	blockApi, _ := api.NewBlockApi(blockRepository, publisherId, blockPool)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi.AddBlockToPool(test.input.block)
	}
}
