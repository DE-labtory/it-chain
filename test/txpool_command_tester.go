package main

import (
	"github.com/it-chain/engine/common/amqp/pubsub"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

func main() {

	config := conf.GetConfiguration()
	client := pubsub.Connect(config.Engine.Amqp)
	defer client.Close()

	txCreateCommand := txpool.TxCreateCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		ICodeID: "1",
		Jsonrpc: "2.0",
		Method:  "invoke",
		Params: txpool.Param{
			Args:     []string{},
			Function: "initA",
		},
	}

	err := client.Publish("Command", "transaction.create", txCreateCommand)

	if err != nil {
		panic(err)
	}
}
