package main

import (
	"encoding/json"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/it-chain-Engine/messaging"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/it-chain/it-chain-Engine/messaging/topic"
)

type EventPublisher struct {
	messaging *messaging.Rabbitmq
}

func NewEventPublisher(messaging *messaging.Rabbitmq) *EventPublisher {
	return &EventPublisher{
		messaging: messaging,
	}
}

func (ep EventPublisher) PublishNewConnEvent(connection conn.Connection) error {

	connInfo := connection.GetConnInfo()

	newConnEvent := event.NewConnEvent{}
	newConnEvent.Id = string(connInfo.Id)
	newConnEvent.Address = connInfo.Address.IP

	b, err := json.Marshal(newConnEvent)

	if err != nil {
		return err
	}

	err = ep.messaging.Publish(topic.MessageCreated.String(), b)

	if err != nil {
		return err
	}

	return nil
}
