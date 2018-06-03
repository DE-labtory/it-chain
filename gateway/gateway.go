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

	//create connection store by bifrost which is it-chain's own lib for implementing p2p network
	ConnectionStore := bifrost.NewConnectionStore()

	//load key
	pri, pub := loadKeyPair(config.Authentication.KeyPath)

	//create amqp Handler
	connectionHandler := NewConnectionCommandHandler(ConnectionStore, pri, pub, rabbitmqClient) // message handler와 구별하기 위해 connection handler로 rename
	messageHandler := NewMessageCommandHandler(ConnectionStore, rabbitmqClient)

	//create gRPC server
	server := NewServer(rabbitmqClient, ConnectionStore, pri, pub)

	// Subscribe amqp server
	// midgard를 사용하여 새 노드 연결 관련 이벤트 구독
	// connectionHandler가 갖는 모든 함수를 실행.
	err := rabbitmqClient.Subscribe("Command", "Connection", connectionHandler)

	if err != nil {
		panic(err)
	}

	//메세지 관련 이벤트 구독
	err = rabbitmqClient.Subscribe("Command", "Messasge", messageHandler)

	if err != nil {
		panic(err)
	}

	// config의 config.yaml에 설정된 grpc gateway의 ip를 서버로 설정한다.
	//bifrost 의 listen 호출을 통해 gRPC 서버를 동작시킨다.
	// 추후 다른 노드에서 실행하는 경우 해당 부분의 ip를 해당 pc의 ip로 바꾸어 주어야 한다.
	server.Listen(config.GrpcGateway.Ip)

	return nil
}

func Stop() {

}
