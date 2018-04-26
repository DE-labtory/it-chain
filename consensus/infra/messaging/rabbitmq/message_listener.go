package rabbitmq

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/consensus/api"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/parliament"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/streadway/amqp"
)

type MessageConsumer struct {
	consensusApi api.ConsensusApi
	messageApi   api.MessageApi
}

func (mc MessageConsumer) HandleStartConsensusEvent(amqpMessage amqp.Delivery) {

	eventMessage := &event.ConsensusStartEvent{}
	err := json.Unmarshal(amqpMessage.Body, &eventMessage)

	if err != nil {
		//error
	}

	err = mc.consensusApi.StartConsensus(parliament.PeerID{eventMessage.UserID}, eventMessage.Block)

	if err != nil {
		//error
	}
}

func (mc MessageConsumer) HandleReceviedConsensusMessageEvent(amqpMessage amqp.Delivery) {

	consensusMsg := &event.ConsensusMessageArriveEvent{}
	err := json.Unmarshal(amqpMessage.Body, &consensusMsg)

	if err != nil {
		//error
	}

	switch consensusMsg.MessageType {
	case event.PREPREPARE:

		preprepareMsg := msg.PreprepareMsg{}
		err := json.Unmarshal(amqpMessage.Body, &preprepareMsg)

		if err != nil {

		}

		mc.consensusApi.ReceivePreprepareMsg(preprepareMsg)
		break

	case event.PREPARE:

		prepareMsg := msg.PrepareMsg{}
		err := json.Unmarshal(amqpMessage.Body, &prepareMsg)

		if err != nil {

		}

		mc.consensusApi.ReceivePrepareMsg(prepareMsg)
		break

	case event.COMMIT:

		commitMsg := msg.CommitMsg{}
		err := json.Unmarshal(amqpMessage.Body, &commitMsg)

		if err != nil {

		}

		mc.consensusApi.ReceiveCommitMsg(commitMsg)
		break
	}
}
