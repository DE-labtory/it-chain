package msg

import (
	"github.com/golang/protobuf/proto"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/protos"
)

type PreprepareMsg struct {
	Consensus consensus.Consensus
	SenderID  string
}

func (pm PreprepareMsg) ToByte() ([]byte, error) {
	data, err := common.Serialize(pm)

	if err != nil {
		return nil, err
	}

	streamMsg := &protos.StreamMsg{}
	streamMsg.Content = &protos.StreamMsg_PreprepareMessage{
		PreprepareMessage: &protos.PreprepareMessage{Data: data}}

	streamData, err := proto.Marshal(streamMsg)

	if err != nil {
		return nil, err
	}

	return streamData, nil
}
