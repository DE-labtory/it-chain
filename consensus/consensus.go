package consensus

import (
	"github.com/it-chain/it-chain-Engine/consensus/infra/messaging/rabbitmq"
	"github.com/it-chain/it-chain-Engine/messaging"
)

func Init() {

	mc := rabbitmq.MessageConsumer{}

	message := messaging.NewMessaging("amqp://guest:guest@localhost:5672/")
	message.Start()

	if err := message.Consume("ConsensusStartEvent", mc.HandleStartConsensusEvent); err != nil {
		panic(err.Error())
	}
}
