package msg

import (
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/protos"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/golang/protobuf/proto"
)

type CommitMsg struct {
	ConsensusID  consensus.ConsensusID
}

func (c CommitMsg) ToByte() ([]byte,error){
	data, err := common.Serialize(c)

	if err != nil{
		return nil, err
	}

	streamMsg := &protos.StreamMsg{}
	streamMsg.Content = &protos.StreamMsg_CommitMessage{
		CommitMessage:&protos.CommitMessage{Data:data}}

	streamData,err := proto.Marshal(streamMsg)

	if err != nil{
		return nil, err
	}

	return streamData, nil
}