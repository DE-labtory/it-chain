package main

import (
	"encoding/json"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/it-chain-Engine/messaging"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/it-chain/it-chain-Engine/messaging/topic"
)

type EventPublisher struct {
	messaging *messaging.Messaging
}

func (ep EventPublisher) PublishConnectionCreatedEvent(connection conn.Connection) error {

	connInfo := connection.GetConnInfo()

	connCreatedEvent := event.ConnectionCreated{}
	connCreatedEvent.Id = string(connInfo.Id)
	connCreatedEvent.Address = connInfo.Address.IP

	b, err := json.Marshal(connCreatedEvent)

	if err != nil {
		return err
	}

	err = ep.messaging.Publish(topic.MessageCreated.String(), b)

	if err != nil {
		return err
	}

	return nil
}
