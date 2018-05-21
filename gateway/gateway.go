package gateway

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/midgard/bus/rabbitmq"
)

func Start() error {

	config := conf.GetConfiguration()

	//create rabbitmq client
	rabbitmqClient := rabbitmq.Connect(config.Common.Messaging.Url)

	//create connection store
	ConnectionStore := bifrost.NewConnectionStore()

	//load key
	pri, pub := loadKeyPair(config.Authentication.KeyPath)

	//createHandler
	commandHandler := NewConnectionCommandHandler(ConnectionStore, pri, pub, rabbitmqClient)
	messageHandler := NewMessageCommandHandler(ConnectionStore, rabbitmqClient)

	//create server
	server := NewServer(rabbitmqClient, ConnectionStore, pri, pub)

	err := rabbitmqClient.Subscribe("Command", "Connection", commandHandler)

	if err != nil {
		panic(err)
	}

	err = rabbitmqClient.Subscribe("Command", "Messasge", messageHandler)

	if err != nil {
		panic(err)
	}

	server.Listen(config.GrpcGateway.Ip)

	return nil
}

func Stop() {

}
