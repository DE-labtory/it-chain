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

	blockRepository := mock.BlockQueryApi{}
	publisherId := "zf"
	eventRepository := mock.EventRepository{}
	eventRepository.LoadFunc = func(aggregate midgard.Aggregate, aggregateID string) error {
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
	blockQueryApi := mock.BlockQueryApi{}
	blockQueryApi.GetLastBlockFunc = func() (blockchain.Block, error) {
		return &blockchain.DefaultBlock{
			Height: blockchain.BlockHeight(12),
		}, nil
	}
	publisherId := "zf"
	eventRepository := mock.EventRepository{}
	eventRepository.LoadFunc = func(aggregate midgard.Aggregate, aggregateID string) error {
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

func TestBlockApi_SyncedCheck(t *testing.T) {
	// TODO:
}

func TestBlockApi_SyncIsProgressing(t *testing.T) {
	// when
	blockQueryApi := mock.BlockQueryApi{}
	eventRepository := mock.EventRepository{}
	eventRepository.LoadFunc = func(aggregate midgard.Aggregate, aggregateID string) error {
		assert.Equal(t, blockchain.BC_SYNC_STATE_AID, aggregateID)
		return nil
	}
	publisherId := "zf"

	// when
	blockApi, _ := api.NewBlockApi(blockQueryApi, eventRepository, publisherId)

	// then
	state := blockApi.SyncIsProgressing()
	assert.Equal(t, blockchain.DONE, state)
}