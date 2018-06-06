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
	// todo change node repo and leader repo after managing peerTable
	nodeRepository := leveldb.NewNodeRepository("path1")
	leaderRepository := leveldb.NewLeaderRepository("path2")
	publisher := rabbitmq.Connect("")
	messageDispatcher := messaging.NewMessageDispatcher(publisher)

	//create amqp Handler
	eventHandler := messaging.NewNodeEventHandler(nodeRepository, leaderRepository, messageDispatcher)
	grpcMessageHandler := messaging.NewGrpcMessageHandler(nodeRepository, leaderRepository, messageDispatcher)

	// Subscribe amqp server
	err1 := rabbitmqClient.Subscribe("Command", "Connection", eventHandler)

	if err1 != nil {
		panic(err1)
	}


	err2 := rabbitmqClient.Subscribe("Command", "message.receive", grpcMessageHandler)

	if err2 != nil {
		panic(err2)
	}
}
