package api_test

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"testing"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/midgard"
	"fmt"
	"errors"
	"github.com/it-chain/it-chain-Engine/blockchain/test/mock"
)

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

	blockRepository := mock.MockBlockQueryApi{}
	publisherId := "zf"
	eventRepository := mock.MockEventRepository{}
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
			height blockchain.BlockHeight
		}
		err error
	} {
		"success": {
			input: struct {
				height blockchain.BlockHeight
			}{height: blockchain.BlockHeight(12),},
			err: nil,
		},
		"block nil test": {
			input: struct {
				height blockchain.BlockHeight
			}{height: blockchain.BlockHeight(144),},
			err: api.ErrNilBlock,
		},
	}
	// When
	blockQueryApi := mock.MockBlockQueryApi{}
	blockQueryApi.GetLastBlockFunc = func() (blockchain.Block, error) {
		return &blockchain.DefaultBlock{
			Height: blockchain.BlockHeight(12),
		}, nil
	}
	publisherId := "zf"
	eventRepository := mock.MockEventRepository{}
	eventRepository.LoadFunc = func(aggregate midgard.Aggregate) error {
		switch v := aggregate.(type) {
		case blockchain.BlockPool:
			aggregate.(blockchain.BlockPool).Add(&blockchain.DefaultBlock{
				Height: blockchain.BlockHeight(12),
				TxList: []blockchain.Transaction{},
			})
			break

		case blockchain.SyncState:
			aggregate.(blockchain.SyncState).SetProgress(blockchain.DONE)
			break
		default:
			return errors.New(fmt.Sprintf("unhandled type [%s]", v))
		}

		return nil
	}

	// When
	blockApi, _ := api.NewBlockApi(blockQueryApi, eventRepository, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		err := blockApi.CheckAndSaveBlockFromPool(test.input.height)

		// Then
		assert.Equal(t, test.err, err)
	}
}