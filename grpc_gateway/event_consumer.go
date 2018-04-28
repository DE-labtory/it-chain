package main

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/streadway/amqp"
)

//Subscribe event and do some logic
type EventConsumer struct {
	messageDeliver MessageDeliver
}

func (ec EventConsumer) MessageDeliveryEvent(amqpMessage <-chan amqp.Delivery) {

	go func() {
		for message := range amqpMessage {

			MessageDelivery := &event.MessageDeliveryEvent{}
			if err := json.Unmarshal(message.Body, MessageDelivery); err != nil {
				// fail to unmarshal event
				return
			}

			ec.messageDeliver.Deliver(MessageDelivery.Recipients, MessageDelivery.Protocol, MessageDelivery.Body)
		}
	}()
}
