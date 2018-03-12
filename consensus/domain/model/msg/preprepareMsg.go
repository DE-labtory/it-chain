package msg

import (
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/protos"
	"github.com/golang/protobuf/proto"
)

type PreprepareMsg struct {
	Consensus       consensus.Consensus
}

func (pm PreprepareMsg) ToByte() ([]byte,error){
	data, err := common.Serialize(pm)

	if err != nil{
		return nil, err
	}

	streamMsg := &protos.StreamMsg{}
	streamMsg.Content = &protos.StreamMsg_PreprepareMessage{
		PreprepareMessage:&protos.PreprepareMessage{Data:data}}

	streamData,err := proto.Marshal(streamMsg)

	if err != nil{
		return nil, err
	}

	return streamData, nil
}