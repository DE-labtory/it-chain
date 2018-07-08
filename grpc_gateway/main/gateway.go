package main

<<<<<<< HEAD:grpc_gateway/main/gateway.go
import (
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/grpc_gateway"
	"github.com/it-chain/it-chain-Engine/grpc_gateway/api"
	"github.com/it-chain/it-chain-Engine/grpc_gateway/infra"
	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/it-chain/midgard/store/leveldb"
)

=======
>>>>>>> b466d51ba979029397eda23887bf103d2b803646:gateway/main/gateway.go
var quit = make(chan bool)

// todo bifrost server kill방법
func Start(ampqUrl string, grpcUrl string, keyPath string) error {

<<<<<<< HEAD:grpc_gateway/main/gateway.go
	//create rabbitmq client
	rabbitmqClient := rabbitmq.Connect(ampqUrl)

	//load key
	pri, pub := infra.LoadKeyPair(keyPath, conf.GetConfiguration().Authentication.KeyType)
	//create gRPC server
	hostService := infra.NewGrpcHostService(pri, pub, rabbitmqClient.Publish)

	//midgard EventStore
	repository := midgard.NewRepo(leveldb.NewEventStore(
		".grpc_gateway/eventStore",
		leveldb.NewSerializer(grpc_gateway.ConnectionCreatedEvent{}, grpc_gateway.ConnectionDisconnectedEvent{}, grpc_gateway.ErrorCreatedEvent{}),
	), rabbitmqClient)

	//createHandler
	connectionApi := api.NewConnectionApi(*repository, hostService) // message handler와 구별하기 위해 connection handler로 rename
	messageApi := api.NewMessageApi(hostService)

	// Subscribe amqp server
	// midgard를 사용하여 새 노드 연결 관련 이벤트 구독
	if err := rabbitmqClient.Subscribe("Command", "connection.*", connectionApi); err != nil {
		panic(err)
	}

	if err := rabbitmqClient.Subscribe("Command", "message.*", messageApi); err != nil {
		panic(err)
	}

	//shutdown grpc_gateway
	go func() {
		for {
			select {
			case <-quit:
				hostService.Stop()
				rabbitmqClient.Close()
				return
			default:
				// Do other stuff
			}
		}
	}()

	// config의 config.yaml에 설정된 grpc gateway의 ip를 서버로 설정한다.
	//bifrost 의 listen 호출을 통해 gRPC 서버를 동작시킨다.
	// 추후 다른 노드에서 실행하는 경우 해당 부분의 ip를 해당 pc의 ip로 바꾸어 주어야 한다.
	hostService.Listen(grpcUrl)

=======
	////create rabbitmq client
	//rabbitmqClient := rabbitmq.Connect(ampqUrl)
	//
	////load key
	//pri, pub := infra.LoadKeyPair(keyPath, conf.GetConfiguration().Authentication.KeyType)
	////create gRPC server
	//hostService := infra.NewGrpcHostService(pri, pub, rabbitmqClient.Publish)
	//
	////midgard EventStore
	//repository := midgard.NewRepo(leveldb.NewEventStore(
	//	".gateway/eventStore",
	//	leveldb.NewSerializer(gateway.ConnectionCreatedEvent{}, gateway.ConnectionDisconnectedEvent{}, gateway.ErrorCreatedEvent{}),
	//), rabbitmqClient)
	//
	////createHandler
	//connectionApi := api.NewConnectionApi(*repository, hostService) // message handler와 구별하기 위해 connection handler로 rename
	//messageApi := api.NewMessageApi(hostService)
	//
	//// Subscribe amqp server
	//// midgard를 사용하여 새 노드 연결 관련 이벤트 구독
	//if err := rabbitmqClient.Subscribe("Command", "connection.*", connectionApi); err != nil {
	//	panic(err)
	//}
	//
	//if err := rabbitmqClient.Subscribe("Command", "message.*", messageApi); err != nil {
	//	panic(err)
	//}
	//
	////shutdown gateway
	//go func() {
	//	for {
	//		select {
	//		case <-quit:
	//			hostService.Stop()
	//			rabbitmqClient.Close()
	//			return
	//		default:
	//			// Do other stuff
	//		}
	//	}
	//}()
	//
	//// config의 config.yaml에 설정된 grpc gateway의 ip를 서버로 설정한다.
	////bifrost 의 listen 호출을 통해 gRPC 서버를 동작시킨다.
	//// 추후 다른 노드에서 실행하는 경우 해당 부분의 ip를 해당 pc의 ip로 바꾸어 주어야 한다.
	//hostService.Listen(grpcUrl)
	//
>>>>>>> b466d51ba979029397eda23887bf103d2b803646:gateway/main/gateway.go
	return nil
}

func Stop() {
	quit <- true
}
