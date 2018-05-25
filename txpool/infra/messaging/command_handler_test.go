package messaging_test

import (
	"testing"

	"os"

	"log"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/infra/messaging"
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

	os.RemoveAll(path)

	repo := midgard.NewRepo(store, client)

	txCommandHandler := messaging.NewTxCommandHandler(*repo, "1")

	//when
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
	}

	//when
	txCommandHandler.HandleTxCreate(txCreatedCommand)

	//then
	tx := &txpool.Transaction{}
	err := repo.Load(tx, "1")
	assert.NoError(t, err)

	log.Println(tx)
}
