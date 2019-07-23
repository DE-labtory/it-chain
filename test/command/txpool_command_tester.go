package main

import (
	"log"

	"github.com/DE-labtory/it-chain/common/command"
	"github.com/DE-labtory/it-chain/common/rabbitmq/rpc"
	"github.com/DE-labtory/it-chain/conf"
	"github.com/DE-labtory/it-chain/txpool"
	"github.com/rs/xid"
)

func main() {

	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	defer client.Close()

	txCreateCommand := command.CreateTransaction{
		TransactionId: xid.New().String(),
		ICodeID:       "bdeshe0e2r74d1hr8pv0",
		Jsonrpc:       "2.0",
		Method:        "invoke",
		Args:          []string{},
		Function:      "initA",
	}

	err := client.Call("transaction.create", txCreateCommand, func(transaction txpool.Transaction, err rpc.Error) {
		log.Printf("created transaction id [%s]", transaction.ID)
	})

	if err != nil {
		log.Println(err)
	}
}
