package adapter_test

import (
	"testing"

	"encoding/hex"

	"reflect"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestBlockResultCommandHandler_HandleBlockResultCommand(t *testing.T) {
	tests := map[string]struct {
		input struct {
			command command.ReturnBlockResult
		}
		err rpc.Error
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
			err: rpc.Error{Message: adapter.ErrBlockIdNil.Error()},
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
			err: rpc.Error{Message: adapter.ErrTxResultsLengthOfZero.Error()},
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
			err: rpc.Error{Message: adapter.ErrTxResultsFail.Error()},
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
			err: rpc.Error{},
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

		value, err := handler.HandleBlockResultCommand(test.input.command)

		assert.Equal(t, err, test.err)
		assert.True(t, reflect.DeepEqual(value, struct{}{}))
	}
}
