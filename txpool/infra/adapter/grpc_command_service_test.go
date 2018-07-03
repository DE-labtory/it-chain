package adapter_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/infra/adapter"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/it-chain-Engine/common"
)

func TestGrpcCommandService_SendLeaderTransactions(t *testing.T) {
	tests := map[string] struct {
		input struct {
			transactions []*txpool.Transaction
			leader txpool.Leader
		}
		err error
	} {
		"success": {
			input: struct {
				transactions []*txpool.Transaction
				leader txpool.Leader
			} {
				transactions: []*txpool.Transaction{&txpool.Transaction{TxId: txpool.TransactionId("zf")}},
				leader: txpool.Leader{LeaderId: txpool.LeaderId{Id: "zf2"}},
			},
			err: nil,
		},
		"transaction empty test": {
			input: struct {
				transactions []*txpool.Transaction
				leader txpool.Leader
			} {
				transactions: []*txpool.Transaction{},
				leader: txpool.Leader{LeaderId: txpool.LeaderId{Id: "zf2"}},
			},
			err: adapter.ErrTxEmpty,
		},
	}

	publisher := func(exchange string, topic string, data interface{}) (err error) {
		txList := &[]*txpool.Transaction{}
		command :=  data.(txpool.GrpcDeliverCommand)

		common.Deserialize(command.Body, txList)

		assert.Equal(t, exchange, "Command")
		assert.Equal(t, topic, "message.deliver")
		assert.Equal(t, 1, len(*txList))

		return nil
	}
	grpcCommandHandler := adapter.NewGrpcCommandService(publisher)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := grpcCommandHandler.SendLeaderTransactions(test.input.transactions, test.input.leader)

		assert.Equal(t, test.err, err)
	}
}
