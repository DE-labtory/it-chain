package gateway

import (
	"os"

	"log"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
)

var (
	ConnectionStore *bifrost.ConnectionStore
	pri             key.PriKey
	pub             key.PubKey
	config          *conf.Configuration
	mq              *rabbitmq.MessageQueue
	s               *Server
	consumer        *Consumer
)

func init() {

	config = conf.GetConfiguration()
	mq = rabbitmq.Connect(config.Common.Messaging.Url)

	ConnectionStore = bifrost.NewConnectionStore()
	pri, pub = loadKeyPair(config.Authentication.KeyPath)
	//mq = rabbitmq.Connect(config.Common.Messaging.Url)

	//publisher
	publisher := NewEventPublisher(mq)

	//amqp event consumer
	consumer = NewAMQPConsumer(ConnectionStore, publisher, pri, pub)
	s = NewServer(consumer, publisher, ConnectionStore, pri, pub)
}

func Start() error {

	mq.Consume(topic.MessageDeliverEvent.String(), consumer.HandleMessageDeliverEvent)
	mq.Consume(topic.ConnCreateCmd.String(), consumer.HandleConnCreateCmd)

	s.Listen(config.GrpcGateway.Ip)

	return nil
}

func Stop() {
	//fmt.Print(mq)

	if mq != nil {
		mq.Close()
	}

	if s != nil {
		s.Stop()
	}

	log.Printf("gateway is closing")
	os.Exit(1)
}
