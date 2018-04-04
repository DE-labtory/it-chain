package rabbitmq

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/consensus/api"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/parliament"
	"github.com/it-chain/it-chain-Engine/messaging/event_message"
	"github.com/streadway/amqp"
)

type MessageConsumer struct {
	consensusApi api.ConsensusApi
	messageApi   api.MessageApi
}

func (mc MessageConsumer) ListenStartConsensusEvent(amqpMessage <-chan amqp.Delivery) {

	go func() {
		for message := range amqpMessage {

			eventMessage := &event_message.StartConsensusEvent{}
			err := json.Unmarshal(message.Body, &eventMessage)

			if err != nil {
				//error
			}

			err = mc.consensusApi.StartConsensus(parliament.PeerID{eventMessage.userID}, eventMessage.block)

			if err != nil {
				//error
			}
		}
	}()
}

func (mc MessageConsumer) ListenReceviedConsensusMessageEvent(amqpMessage <-chan amqp.Delivery) {

	go func() {
		for message := range amqpMessage {

			consensusMsg := &event_message.ReceviedConsensusMessageEvent{}
			err := json.Unmarshal(message.Body, &consensusMsg)

			if err != nil {
				//error
			}

			switch consensusMsg.MessageType {
			case event_message.PREPREPARE:

				preprepareMsg := msg.PreprepareMsg{}
				err := json.Unmarshal(message.Body, &preprepareMsg)

				if err != nil {

				}

				mc.consensusApi.ReceivePreprepareMsg(preprepareMsg)
				break

			case event_message.PREPARE:

				prepareMsg := msg.PrepareMsg{}
				err := json.Unmarshal(message.Body, &prepareMsg)

				if err != nil {

				}

				mc.consensusApi.ReceivePrepareMsg(prepareMsg)
				break

			case event_message.COMMIT:

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
