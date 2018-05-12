//Subscribe event and do corresponding logic

package gateway

import (
	"encoding/json"
	"log"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
	"github.com/streadway/amqp"
)

type Consumer struct {
	ConnectionStore *bifrost.ConnectionStore
	publisher       *EventPublisher
	priKey          key.PriKey
	pubKey          key.PubKey
}

func NewAMQPConsumer(ConnectionStore *bifrost.ConnectionStore, publisher *EventPublisher, pri key.PriKey, pub key.PubKey) *Consumer {
	return &Consumer{
		ConnectionStore: ConnectionStore,
		publisher:       publisher,
		priKey:          pri,
		pubKey:          pub,
	}
}

func (c Consumer) HandleMessageDeliverEvent(amqpMessage amqp.Delivery) {

	MessageDelivery := &event.MessageDeliverEvent{}
	if err := json.Unmarshal(amqpMessage.Body, MessageDelivery); err != nil {
		// fail to unmarshal event
		return
	}

	for _, recipient := range MessageDelivery.Recipients {
		connection := c.ConnectionStore.GetConnection(bifrost.ConnID(recipient))

		if connection != nil {
			connection.Send(MessageDelivery.Body, MessageDelivery.Protocol, nil, nil)
		}
	}
}

func (c Consumer) ServeRequest(msg bifrost.Message) {

}

func (c Consumer) ServeError(conn bifrost.Connection, err error) {

}

func (c Consumer) HandleConnCreateCmd(amqpMessage amqp.Delivery) {

	log.Println("ConnCreatedCmd")
	ConnCreateCmd := &event.ConnCreateCmd{}

	if err := json.Unmarshal(amqpMessage.Body, ConnCreateCmd); err != nil {
		c.publisher.PublishGatewayErrorEvent(topic.ConnCreateCmd.String(), err.Error())
		return
	}

	clientOpt := client.ClientOpts{
		Ip:     ConnCreateCmd.Address,
		PriKey: c.priKey,
		PubKey: c.pubKey,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	connection, err := client.Dial(ConnCreateCmd.Address, clientOpt, grpcOpt)

	if err != nil {
		c.publisher.PublishGatewayErrorEvent(topic.ConnCreateCmd.String(), err.Error())
		return
	}

	connection.Handle(c)
	c.ConnectionStore.AddConnection(connection)

	c.publisher.PublishConnCreatedEvent(connection)

	go func() {
		defer connection.Close()

		if err := connection.Start(); err != nil {
			connection.Close()
		}
		log.Printf("connection are deleted")
	}()
}
