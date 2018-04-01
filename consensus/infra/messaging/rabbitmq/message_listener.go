package rabbitmq

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/consensus/api"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/parliament"
	"github.com/streadway/amqp"
)

type ConsensusMessageType int

const (
	PREPREPARE ConsensusMessageType = 0
	PREPARE    ConsensusMessageType = 1
	COMMIT     ConsensusMessageType = 2
)

type StartConsensusEvent struct {
	block  []byte
	userID string
}

type ReceviedConsensusMessageEvent struct {
	messageType ConsensusMessageType
	messageBody []byte
}

type MessageConsumer struct {
	consensusApi api.ConsensusApi
	messageApi   api.MessageApi
}

func (mc MessageConsumer) ListenStartConsensusEvent(amqpMessage <-chan amqp.Delivery) {

	go func() {
		for message := range amqpMessage {

			eventMessage := &StartConsensusEvent{}
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

			eventMessage := &ReceviedConsensusMessageEvent{}
			err := json.Unmarshal(message.Body, &eventMessage)

			if err != nil {
				//error
			}

			switch eventMessage.messageType {
			case PREPREPARE:

				preprepareMsg := msg.PreprepareMsg{}
				err := json.Unmarshal(message.Body, &preprepareMsg)

				if err != nil {

				}

				mc.consensusApi.ReceivePreprepareMsg(preprepareMsg)
				break

			case PREPARE:

				prepareMsg := msg.PrepareMsg{}
				err := json.Unmarshal(message.Body, &prepareMsg)

				if err != nil {

				}

				mc.consensusApi.ReceivePrepareMsg(prepareMsg)
				break

			case COMMIT:

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
