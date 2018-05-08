package main

import (
	"encoding/json"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/streadway/amqp"
)

//Subscribe event and do corresponding logic
type EventConsumer struct {
	connectionStore *bifrost.ConnectionStore
}

func NewEventConsumer(connectionStore *bifrost.ConnectionStore) *EventConsumer {
	return &EventConsumer{
		connectionStore: connectionStore,
	}
}

func (ec EventConsumer) HandleMessageDeliverEvent(amqpMessage amqp.Delivery) {

	MessageDelivery := &event.MessageDeliverEvent{}
	if err := json.Unmarshal(amqpMessage.Body, MessageDelivery); err != nil {
		// fail to unmarshal event
		return
	}

	ec.deliver(MessageDelivery.Recipients, MessageDelivery.Protocol, MessageDelivery.Body)
}

func (ec EventConsumer) deliver(recipients []string, protocol string, data []byte) error {

	for _, recipient := range recipients {
		connection := ec.connectionStore.GetConnection(bifrost.ConnID(recipient))

		if connection != nil {
			connection.Send(data, protocol, nil, nil)
		}
	}

	return nil
}
