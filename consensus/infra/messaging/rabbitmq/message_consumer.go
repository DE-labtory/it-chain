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
}

func (mc MessageConsumer) ListenStartConsensusEvent(amqpMessage <-chan amqp.Delivery) {

	go func() {
		for message := range amqpMessage {

			eventMessage := &event.ConsensusStartEvent{}
			err := json.Unmarshal(message.Body, &eventMessage)

			if err != nil {
				//error
			}

			err = mc.consensusApi.StartConsensus(parliament.PeerID{eventMessage.UserID}, eventMessage.Block)

			if err != nil {
				//error
			}
		}
	}()
}

func (mc MessageConsumer) ListenReceviedConsensusMessageEvent(amqpMessage <-chan amqp.Delivery) {

	go func() {
		for message := range amqpMessage {

			consensusMsg := &event.ConsensusMessageArriveEvent{}
			err := json.Unmarshal(message.Body, &consensusMsg)

			if err != nil {
				//error
			}

			switch consensusMsg.MessageType {
			case event.PREPREPARE:

				preprepareMsg := msg.PreprepareMsg{}
				err := json.Unmarshal(message.Body, &preprepareMsg)

				if err != nil {

				}

				mc.consensusApi.ReceivePreprepareMsg(preprepareMsg)
				break

			case event.PREPARE:

				prepareMsg := msg.PrepareMsg{}
				err := json.Unmarshal(message.Body, &prepareMsg)

				if err != nil {

				}

				mc.consensusApi.ReceivePrepareMsg(prepareMsg)
				break

			case event.COMMIT:

				commitMsg := msg.CommitMsg{}
				err := json.Unmarshal(message.Body, &commitMsg)

				if err != nil {

				}

				mc.consensusApi.ReceiveCommitMsg(commitMsg)
				break
			}
		}
	}()
}
