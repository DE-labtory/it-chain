package service

import (
	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
)

//함수의 로직이 똑같기 때문에 재사용 하는 방향으로 구현되어야함
func CheckPreparePolicy(consensus cs.Consensus, msgPool msg.MsgPool) bool {
	if consensus.IsPrepareState() {
		numberOfRepresentatives := float64(len(consensus.Representatives))
		//
		prepareMsgs := msgPool.FindPrepareMsgsByConsensusID(consensus.ConsensusID)
		numberOfPrepareMsgs := len(prepareMsgs)

		if numberOfPrepareMsgs > int((numberOfRepresentatives)/3)+1 {
			return true
		}
		return false
	}

	return false
}

func CheckCommitPolicy(consensus cs.Consensus, msgPool msg.MsgPool) bool {
	if consensus.IsCommitState() {
		numberOfRepresentatives := float64(len(consensus.Representatives))
		//
		commitMsgs := msgPool.FindCommitMsgsByConsensusID(consensus.ConsensusID)
		numberOfPrepareMsgs := len(commitMsgs)

		if numberOfPrepareMsgs > int((numberOfRepresentatives)/3)+1 {
			return true
		}
		return false
	}

	return false
}
