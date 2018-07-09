package api_test

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"testing"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/midgard"
	"fmt"
	"errors"
)

type MockBlockQueryApi struct {
	GetLastBlockFunc func() (blockchain.Block, error)
}
func (br MockBlockQueryApi) GetLastBlock() (blockchain.Block, error) {
	return br.GetLastBlockFunc()
}
func (br MockBlockQueryApi) GetBlockByHeight(blockHeight uint64) (blockchain.Block, error) { return nil, nil }
func (br MockBlockQueryApi) GetBlockBySeal(seal []byte) (blockchain.Block, error) { return nil, nil }
func (br MockBlockQueryApi) GetBlockByTxID(txid string) (blockchain.Block, error) { return nil, nil }
func (br MockBlockQueryApi) GetTransactionByTxID(txid string) (blockchain.Transaction, error) { return nil, nil }

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

	blockRepository := MockBlockQueryApi{}
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
	blockQueryApi := MockBlockQueryApi{}
	blockQueryApi.GetLastBlockFunc = func() (blockchain.Block, error) {
		return &blockchain.DefaultBlock{
			Height: blockchain.BlockHeight(12),
		}, nil
	}
	publisherId := "zf"
	eventRepository := MockEventRepository{}
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