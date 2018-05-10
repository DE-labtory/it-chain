package rabbitmq

import (
	"fmt"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/domain/service"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
)

type Publish func(topic string, data []byte) error

type MessageService struct {
	Publish Publish
}

func NewRabbitmqMessageService(publish Publish) *MessageService {
	return &MessageService{Publish: publish}
}

func (mApi *MessageService) BroadCastMsg(Msg service.Serializable, representatives []*consensus.Representative) {

	data, err := Msg.ToByte()

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	IDs := make([]string, 0)

	for _, representative := range representatives {
		IDs = append(IDs, representative.GetIdString())
	}

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	serializedData, err := common.Serialize(event.ConsensusMessagePublishEvent{Ids: IDs, Data: data})

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	err = mApi.Publish(topic.ConsensusMessagePublishEvent.String(), serializedData)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
}

func (mApi *MessageService) ConfirmedBlock(block consensus.Block) {

	err := mApi.Publish(topic.BlockConfirmEvent.String(), block)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
}
