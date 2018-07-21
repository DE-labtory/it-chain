package adapter_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

func TestBlockProposeCommandHandler_HandleProposeBlockCommand(t *testing.T) {
	tests := map[string]struct {
		input struct {
			command command.ProposeBlock
		}
		err error
	}{
		"command with emtpy transactions test": {
			input: struct {
				command command.ProposeBlock
			}{
				command: command.ProposeBlock{
					CommandModel: midgard.CommandModel{ID: "111"},
					TxList:       nil,
				},
			},
			err: adapter.ErrCommandTransactions,
		},
		"transactions which have length of 0 test": {
			input: struct {
				command command.ProposeBlock
			}{
				command: command.ProposeBlock{
					CommandModel: midgard.CommandModel{ID: "111"},
					TxList:       make([]command.Tx, 0),
				},
			},
			err: adapter.ErrCommandTransactions,
		},
		"transactions which have missing properties test": {
			input: struct {
				command command.ProposeBlock
			}{
				command: command.ProposeBlock{
					CommandModel: midgard.CommandModel{ID: "111"},
					TxList: []command.Tx{
						command.Tx{ID: "", PeerID: ""},
					},
				},
			},
			err: adapter.ErrTxHasMissingProperties,
		},
		"successfully pass txlist to block api": {
			input: struct {
				command command.ProposeBlock
			}{
				command: command.ProposeBlock{
					CommandModel: midgard.CommandModel{ID: "111"},
					TxList: []command.Tx{
						command.Tx{
							ID:        "1",
							Status:    1,
							PeerID:    "2",
							TimeStamp: time.Now(),
							Jsonrpc:   "123",
							Method:    "invoke",
							Function:  "function1",
							Args:      []string{"arg1", "arg2"},
							Signature: []byte{0x1},
						},
					},
				},
			},
			err: nil,
		},
	}

	blockApi := mock.BlockApi{}
	blockApi.CreateBlockFunc = func(txList []blockchain.Transaction) error {
		tx := txList[0]
		txContentBytes, _ := tx.GetContent()
		content := struct {
			ID        string
			Status    blockchain.Status
			PeerID    string
			Timestamp time.Time
			TxData    *blockchain.TxData
		}{}
		json.Unmarshal(txContentBytes, &content)

		// then
		assert.Equal(t, "1", tx.GetID())
		assert.Equal(t, "2", content.PeerID)
		assert.Equal(t, blockchain.Status(1), content.Status)
		assert.Equal(t, "123", content.TxData.Jsonrpc)
		assert.Equal(t, blockchain.Invoke, content.TxData.Method)
		assert.Equal(t, "function1", content.TxData.Params.Function)
		assert.Equal(t, []string{"arg1", "arg2"}, content.TxData.Params.Args)

		return nil
	}

	commandHandler := adapter.NewBlockProposeCommandHandler(blockApi, "solo")

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := commandHandler.HandleProposeBlockCommand(test.input.command)

		assert.Equal(t, err, test.err)
	}
}
