package gateway

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/midgard/bus/rabbitmq"
)

var quit = make(chan bool)

// todo bifrost server kill방법
func Start(ampqUrl string, grpcUrl string, keyPath string) error {

	//create rabbitmq client
	rabbitmqClient := rabbitmq.Connect(ampqUrl)

	//create connection store by bifrost which is it-chain's own lib for implementing p2p network
	ConnectionStore := bifrost.NewConnectionStore()

	pri, pub := LoadKeyPair(keyPath)
	//load key

	//createHandler
	connectionHandler := NewConnectionCommandHandler(ConnectionStore, pri, pub, rabbitmqClient) // message handler와 구별하기 위해 connection handler로 rename
	messageHandler := NewMessageCommandHandler(ConnectionStore, rabbitmqClient)

	//create gRPC server
	server := NewServer(rabbitmqClient, ConnectionStore, pri, pub)

	// Subscribe amqp server
	// midgard를 사용하여 새 노드 연결 관련 이벤트 구독
	err := rabbitmqClient.Subscribe("Command", "Connection", connectionHandler)

	if err != nil {
		panic(err)
	}

	//메세지 관련 이벤트 구독
	err = rabbitmqClient.Subscribe("Command", "Messasge", messageHandler)

	if err != nil {
		panic(err)
	}

	//shutdown gateway
	go func() {
		for {
			select {
			case <-quit:
				server.Stop()
				rabbitmqClient.Close()
				return
			default:
				// Do other stuff
			}
		}
	}()

	// config의 config.yaml에 설정된 grpc gateway의 ip를 서버로 설정한다.
	// 추후 다른 노드에서 실행하는 경우 해당 부분의 ip를 해당 pc의 ip로 바꾸어 주어야 한다.
	server.Listen(grpcUrl)

	return nil
}

func Stop() {
	quit <- true
}
