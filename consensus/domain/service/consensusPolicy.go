package service

import (
	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
)

func CheckPreparePolicy(consensus cs.Consensus,msgPool msg.MsgPool) bool{
	//NumberOfRepresentatives := len(consensus.Representatives)
	//
	//msgPool.FindPrepareMsgsByConsensusID()
	return true
}