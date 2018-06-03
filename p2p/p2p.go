package p2p

import (
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/it-chain/it-chain-Engine/p2p/infra/messaging"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
)

func init() {
	//myIp := conf.GetConfiguration().Common.NodeIp

	config := conf.GetConfiguration()
	//create rabbitmq client
	rabbitmqClient := rabbitmq.Connect(config.Common.Messaging.Url)

	nodeRepository := leveldb.NewNodeRepository("path")
	leaderRepository := leveldb.NewLeaderRepository("path")
	publisher := rabbitmq.Connect("")
	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	//create amqp Handler
	connectionHandler := messaging.NewNodeEventHandler(nodeRepository, leaderRepository) // message handler와 구별하기 위해 connection handler로 rename
	messageHandler := NewMessageCommandHandler(ConnectionStore, rabbitmqClient)

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

	//myNode := NewNode(myIp, nodeId) // nodeId는 어디서 할당?
	////messageDispatcher := messaging.NewMessageDispatcher(midgard.Publisher) midgard 주입부분 => midgard doc 완성 후로 보류
	//repository := leveldb.NewNodeRepository("path") //repository 객체 생성
	//leaderSelectionApi := NewLeaderSelectionApi(repository, messageDispatcher, myNode)
	//
	//// 해당 노드를 leveldb에 저장
	//repository.save(myNode)

}
