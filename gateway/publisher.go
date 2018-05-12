package gateway

import (
	"encoding/json"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
)

type Publisher interface {
	Publish(topic string, body []byte) error
}

type EventPublisher struct {
	publisher Publisher
}

func NewEventPublisher(mq Publisher) *EventPublisher {
	return &EventPublisher{
		publisher: mq,
	}
}

func (p EventPublisher) PublishConnCreatedEvent(connection bifrost.Connection) error {

	ConnCreateEvent := event.ConnCreateEvent{
		Id:      string(connection.GetID()),
		Address: connection.GetIP(),
	}

	b, err := json.Marshal(ConnCreateEvent)

	if err != nil {
		return err
	}

	err = p.publisher.Publish(topic.NewConnEvent.String(), b)

	if err != nil {
		return err
	}

	return nil
}

func (p EventPublisher) PublishGatewayErrorEvent(event string, error string) error {

	errorEvent := &ErrorEvent{
		Err:   error,
		Event: event,
	}

	b, err := json.Marshal(errorEvent)

	if err != nil {
		return err
	}

	err = p.publisher.Publish("GatewayErrorEvent", b)

	if err != nil {
		return err
	}

	return nil
}
