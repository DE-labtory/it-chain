package api

import (
	"fmt"

	"github.com/it-chain/it-chain-Engine/common"
	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/messaging"
	"github.com/it-chain/it-chain-Engine/messaging/event"
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

	serializedData, err := common.Serialize(messaging.Sendable{Ids: IDs, Data: data})

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	err = mApi.Publish(event.MessageCreated.String(), serializedData)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
}

func (mApi *MessageApi) ConfirmedBlock(block cs.Block) {

}
