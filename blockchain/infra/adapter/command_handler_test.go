package adapter_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/it-chain/it-chain-Engine/blockchain/test/mock"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestCommandHandler_HandleConfirmBlockCommand(t *testing.T) {
	tests := map[string]struct {
		input struct {
			command blockchain.ConfirmBlockCommand
		}
		err error
	}{
		"success": {
			input: struct {
				command blockchain.ConfirmBlockCommand
			}{
				command: blockchain.ConfirmBlockCommand{
					CommandModel: midgard.CommandModel{ID: "zf"},
					Block: &blockchain.DefaultBlock{
						Height: 99887,
					},
				},
			},
			err: nil,
		},
		"block nil error test": {
			input: struct {
				command blockchain.ConfirmBlockCommand
			}{
				command: blockchain.ConfirmBlockCommand{
					CommandModel: midgard.CommandModel{ID: "zf"},
					Block:        nil,
				},
			},
			err: adapter.ErrBlockNil,
		},
	}

	blockApi := mock.BlockApi{}
	blockApi.StageBlockFunc = func(block blockchain.Block) error {
		assert.Equal(t, block.GetHeight(), uint64(99887))
		return nil
	}

	commandHandler := adapter.NewCommandHandler(blockApi)
	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := commandHandler.HandleConfirmBlockCommand(test.input.command)

		assert.Equal(t, err, test.err)
	}

}
