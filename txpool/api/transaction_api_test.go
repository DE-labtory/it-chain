package api_test

import (
	"testing"

	"os"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/api"
	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/it-chain/midgard/store/leveldb"
	"github.com/stretchr/testify/assert"
)

//need rabbitmq
func TestTxCommandHandler_HandleTxCreate(t *testing.T) {

	//given
	client := rabbitmq.Connect("")
	path := "test"
	store := leveldb.NewEventStore(path, leveldb.NewSerializer(txpool.TxCreatedEvent{}, txpool.TxDeletedEvent{}))

	defer os.RemoveAll(path)

	repo := midgard.NewRepo(store, client)

	txCommandHandler := api.NewTransactionApi(repo, "1")

	//when
	txID := "123"
	txCreatedCommand := txpool.TxCreateCommand{
		TxData: txpool.TxData{
			ID:      "1",
			ICodeID: "123",
			Params: txpool.Param{
				Args:     []string{},
				Function: "func1",
			},
			Method:  "invoke",
			Jsonrpc: "json1.0",
		},
		CommandModel: midgard.CommandModel{
			ID: txID,
		},
	}

	//when
	txCommandHandler.CreateTransaction(txCreatedCommand)

	//then
	tx := &txpool.Transaction{}
	err := repo.Load(tx, txID)
	assert.NoError(t, err)

	assert.Equal(t, tx.TxData, txCreatedCommand.TxData)
	assert.Equal(t, tx.TxId, tx.TxId)
}
