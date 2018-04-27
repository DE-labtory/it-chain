package main

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/streadway/amqp"
)

//Subscribe event and do corresponding logic
type EventConsumer struct {
	messageDeliver MessageDeliver
}

func (ec EventConsumer) HandleMessageDeliverEvent(amqpMessage amqp.Delivery) {

	MessageDelivery := &event.MessageDeliverEvent{}
	if err := json.Unmarshal(amqpMessage.Body, MessageDelivery); err != nil {
		// fail to unmarshal event
		return
	}

	ec.messageDeliver.Deliver(MessageDelivery.Recipients, MessageDelivery.Protocol, MessageDelivery.Body)
}
