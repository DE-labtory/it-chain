package main

import (
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

func main() {

	config := conf.GetConfiguration()
	client := pubsub.NewTopicPublisher(config.Engine.Amqp, "Command")
	defer client.Close()

	txCreateCommand := command.CreateTransaction{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		ICodeID:  "1",
		Jsonrpc:  "2.0",
		Method:   "invoke",
		Args:     []string{},
		Function: "initA",
	}

	err := client.Publish("transaction.create", txCreateCommand)

	if err != nil {
		panic(err)
	}
}
