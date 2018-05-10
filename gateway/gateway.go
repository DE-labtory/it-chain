package gateway

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
)

func Start() error {

	var (
		ConnectionStore *bifrost.ConnectionStore
		pri             key.PriKey
		pub             key.PubKey
		config          *conf.Configuration
		mq              *rabbitmq.MessageQueue
	)

	config = conf.GetConfiguration()
	ConnectionStore = bifrost.NewConnectionStore()
	pri, pub = loadKeyPair(config.Authentication.KeyPath)
	mq = rabbitmq.Connect(config.Common.Messaging.Url)

	//publisher
	publisher := NewAMQPPublisher(mq)

	//muxer
	muxer := NewGatewayMux(publisher)

	//amqp event consumer
	consumer := NewAMQPConsumer(ConnectionStore, muxer, publisher, pri, pub)

	mq.Consume(topic.MessageDeliverEvent.String(), consumer.HandleMessageDeliverEvent)
	mq.Consume(topic.ConnCreateCmd.String(), consumer.HandleConnCreateCmd)

	//server
	server := NewServer(muxer, pri, pub)
	server.Listen(config.GrpcGateway.Ip)

	return nil
}
