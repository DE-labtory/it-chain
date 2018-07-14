package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

type MockBlockQueryApi struct {
	GetLastBlockFunc         func() (blockchain.Block, error)
	GetStagedBlockByIdFunc   func(blockId string) (blockchain.Block, error)
	GetLastCommitedBlockFunc func() (blockchain.Block, error)
}

func (br MockBlockQueryApi) GetLastBlock() (blockchain.Block, error) {
	return br.GetLastBlockFunc()
}
func (br MockBlockQueryApi) GetBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return nil, nil
}
func (br MockBlockQueryApi) GetStagedBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return nil, nil
}
func (br MockBlockQueryApi) GetStagedBlockById(blockId string) (blockchain.Block, error) {
	return br.GetStagedBlockByIdFunc(blockId)
}
func (br MockBlockQueryApi) GetLastCommitedBlock() (blockchain.Block, error) {
	return br.GetLastCommitedBlockFunc()
}

type MockEventRepository struct {
	LoadFunc func(aggregate midgard.Aggregate) error
}

func (er MockEventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return er.LoadFunc(aggregate)
}
func (er MockEventRepository) Save(aggregateID string, events ...midgard.Event) error { return nil }
func (er MockEventRepository) Close()                                                 {}

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

	blockQueryApi := MockBlockQueryApi{}
	publisherId := "zf"

	blockApi, _ := api.NewBlockApi(blockQueryApi, publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi.AddBlockToPool(test.input.block)
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
	blockQueryApi := MockBlockQueryApi{}
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
