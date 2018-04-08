package consensus

import (
	"github.com/it-chain/it-chain-Engine/consensus/infra/messaging/rabbitmq"
	"github.com/it-chain/it-chain-Engine/messaging"
)

func Init() {

	mc := rabbitmq.MessageConsumer{}

	message := messaging.NewMessaging("amqp://guest:guest@localhost:5672/")
	message.Start()

	event, err := message.Consume("ConsensusStartEvent")
	mc.ListenReceviedConsensusMessageEvent(event)

	if err != nil {

	}
}
