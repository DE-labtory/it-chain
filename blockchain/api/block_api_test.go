package api_test

import (
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"testing"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/midgard"
)

type MockBlockRepository struct {}

func (br MockBlockRepository) Close() {}
func (br MockBlockRepository) GetValidator() common.Validator { return nil }
func (br MockBlockRepository) AddBlock(block common.Block) error { return nil }
func (br MockBlockRepository) GetLastBlock(block common.Block) error { return nil }
func (br MockBlockRepository) NewEmptyBlock() (blockchain.Block, error) {return nil, nil}

type MockEventRepository struct {
	LoadFunc func(aggregate midgard.Aggregate) error
}

func (er MockEventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return er.LoadFunc(aggregate)
}
func (er MockEventRepository) Save(aggregateID string, events ...midgard.Event) error {return nil}
func (er MockEventRepository) Close() {}

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
	}

	blockRepository := MockBlockRepository{}
	publisherId := "zf"
	eventRepository := MockEventRepository{}
	eventRepository.LoadFunc = func(aggregate midgard.Aggregate) error {
		// predefine block pool
		aggregate = blockchain.NewBlockPool()
		return nil
	}

	blockApi, _ := api.NewBlockApi(blockRepository, eventRepository, publisherId)

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
				Height: blockchain.BlockHeight(144),
			}},
			err: api.ErrNilBlock,
		},
	}
	// When
	blockRepository := MockBlockRepository{}
	publisherId := "zf"
	eventRepository := MockEventRepository{}
	eventRepository.LoadFunc = func(aggregate midgard.Aggregate) error {
		aggregate.(blockchain.BlockPool).Add(&blockchain.DefaultBlock{
			Height: blockchain.BlockHeight(12),
		})
		return nil
	}

	// When
	blockApi, _ := api.NewBlockApi(blockRepository, eventRepository, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		err := blockApi.CheckAndSaveBlockFromPool(test.input.block)

		// Then
		assert.Equal(t, test.err, err)
	}
}