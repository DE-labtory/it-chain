package gateway

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/midgard/bus/rabbitmq"
)

// todo bifrost server kill방법
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

	//create gRPC server
	server := NewServer(rabbitmqClient, ConnectionStore, pri, pub)

	// Subscribe amqp server
	err := rabbitmqClient.Subscribe("Command", "Connection", commandHandler)

	if err != nil {
		panic(err)
	}

	err = rabbitmqClient.Subscribe("Command", "Messasge", messageHandler)

	if err != nil {
		panic(err)
	}

	// config의 config.yaml에 설정된 grpc gateway의 ip를 서버로 설정한다.
	// 추후 다른 노드에서 실행하는 경우 해당 부분의 ip를 해당 pc의 ip로 바꾸어 주어야 한다.
	server.Listen(config.GrpcGateway.Ip)

	return nil
}

func Stop() {

}
