package p2p

import (
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/it-chain/it-chain-Engine/p2p/infra/messaging"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"src/github.com/it-chain/midgard"
	leveldb2 "src/github.com/it-chain/midgard/store/leveldb"
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


	leaderApi := api.NewLeaderApi(eventRepository, messageDispatcher, myInfo)

	//create amqp Handler
	eventHandler := messaging.NewNodeEventHandler(nodeRepository, leaderRepository)
	grpcMessageHandler := messaging.NewGrpcMessageHandler(leaderApi, nodeApi, messageDispatcher)

	// Subscribe amqp server
	err1 := rabbitmqClient.Subscribe("Command", "connection.*", eventHandler)

	if err1 != nil {
		panic(err1)
	}


	err2 := rabbitmqClient.Subscribe("Command", "message.*", grpcMessageHandler)

	if err2 != nil {
		panic(err2)
	}
}
