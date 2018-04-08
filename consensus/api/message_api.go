package api

import (
	"fmt"

	"github.com/it-chain/it-chain-Engine/common"
	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/it-chain/it-chain-Engine/messaging/topic"
)

type Serializable interface {
	ToByte() ([]byte, error)
}

type Publish func(topic string, data []byte) error

type MessageApi struct {
	Publish Publish
}

func NewMessageApi(publish Publish) *MessageApi {
	return &MessageApi{Publish: publish}
}

func (mApi *MessageApi) BroadCastMsg(Msg Serializable, representatives []*cs.Representative) {

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

func (mApi *MessageApi) ConfirmedBlock(block cs.Block) {

	err := mApi.Publish(topic.BlockConfirmEvent.String(), block)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
}
