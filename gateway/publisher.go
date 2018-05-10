package gateway

import (
	"encoding/json"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
)

type AMQPPublisher struct {
	mq *rabbitmq.MessageQueue
}

func NewAMQPPublisher(mq *rabbitmq.MessageQueue) *AMQPPublisher {
	return &AMQPPublisher{
		mq: mq,
	}
}

func (p AMQPPublisher) ConnCreatedEvent(connection bifrost.Connection) error {

	newConnEvent := event.NewConnEvent{
		Id:      string(connection.GetID()),
		Address: connection.GetIP(),
	}

	b, err := json.Marshal(newConnEvent)

	if err != nil {
		return err
	}

	err = p.mq.Publish(topic.NewConnEvent.String(), b)

	if err != nil {
		return err
	}

	return nil
}
