package adapter_test

import (
	"testing"

	"encoding/hex"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

func TestBlockResultCommandHandler_HandleBlockResultCommand(t *testing.T) {
	tests := map[string]struct {
		input struct {
			command command.ReturnBlockResult
		}
		err error
	}{
		"test when block id nil": {
			input: struct {
				command command.ReturnBlockResult
			}{command.ReturnBlockResult{
				CommandModel: midgard.CommandModel{
					ID: "",
				},
				TxResultList: mock.GetTxResults(),
			}},
			err: adapter.ErrBlockIdNil,
		},
		"test when length of tx results 0": {
			input: struct {
				command command.ReturnBlockResult
			}{command.ReturnBlockResult{
				CommandModel: midgard.CommandModel{
					ID: "block_id1",
				},
				TxResultList: mock.GetZeroLengthTxResults(),
			}},
			err: adapter.ErrTxResultsLengthOfZero,
		},
		"test when one of tx results failed": {
			input: struct {
				command command.ReturnBlockResult
			}{command.ReturnBlockResult{
				CommandModel: midgard.CommandModel{
					ID: "block_id1",
				},
				TxResultList: mock.GetFailedTxResults(),
			}},
			err: adapter.ErrTxResultsFail,
		},
		"successfully commit block": {
			input: struct {
				command command.ReturnBlockResult
			}{command.ReturnBlockResult{
				CommandModel: midgard.CommandModel{
					ID: "block_id1",
				},
				TxResultList: mock.GetTxResults(),
			}},
			err: nil,
		},
	}

	repo := mock.EventRepository{}
	repo.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, events[0].GetID(), hex.EncodeToString([]byte("block_id1")))
		return nil
	}

	repo.CloseFunc = func() {}

	eventstore.InitForMock(repo)
	defer eventstore.Close()

	blockQueryService := mock.BlockQueryService{}
	blockQueryService.GetStagedBlockByIdFunc = func(blockId string) (blockchain.DefaultBlock, error) {
		assert.Equal(t, blockId, "block_id1")
		return mock.GetStagedBlockWithId(blockId), nil
	}

	handler := adapter.NewBlockResultCommandHandler(blockQueryService)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := handler.HandleBlockResultCommand(test.input.command)

		assert.Equal(t, err, test.err)
	}
}
