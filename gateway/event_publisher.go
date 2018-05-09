package main

import (
	"encoding/json"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/it-chain/it-chain-Engine/messaging/topic"
)

func PublishNewConnEvent(connection bifrost.Connection) error {

	newConnEvent := event.NewConnEvent{
		Id:      string(connection.GetID()),
		Address: connection.GetIP(),
	}

	b, err := json.Marshal(newConnEvent)

	if err != nil {
		return err
	}

	err = mq.Publish(topic.NewConnEvent.String(), b)

	if err != nil {
		return err
	}

	return nil
}
