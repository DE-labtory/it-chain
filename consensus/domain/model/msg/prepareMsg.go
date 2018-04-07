package msg

import (
	"github.com/golang/protobuf/proto"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/protos"
)

type PrepareMsg struct {
	ConsensusID consensus.ConsensusID
	Block       consensus.Block
	SenderID    string
}

func (p PrepareMsg) ToByte() ([]byte, error) {
	data, err := common.Serialize(p)

	if err != nil {
		return nil, err
	}

	streamMsg := &protos.StreamMsg{}
	streamMsg.Content = &protos.StreamMsg_PrepareMessage{
		PrepareMessage: &protos.PrepareMessage{Data: data}}

	streamData, err := proto.Marshal(streamMsg)

	if err != nil {
		return nil, err
	}

	return streamData, nil
}
