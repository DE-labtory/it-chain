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
	eventHandler := messaging.NewNodeEventHandler(nodeRepository, leaderRepository, messageDispatcher)

	// Subscribe amqp server
	err := rabbitmqClient.Subscribe("Command", "Connection", eventHandler)

	if err != nil {
		panic(err)
	}
}
