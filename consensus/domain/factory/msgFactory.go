package factory

import (
	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/msg"
)

func CreatePreprepareMsg(consensus cs.Consensus) msg.PreprepareMsg{
	return msg.PreprepareMsg{
		Consensus: consensus,
	}
}

func CreatePrepareMsg(consensus cs.Consensus) msg.PrepareMsg{
	return msg.PrepareMsg{
		Block:consensus.Block,
		ConsensusID: consensus.ConsensusID,
	}
}

func CreateMsgPool() *msg.MsgPool{
	return &msg.MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]msg.PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]msg.CommitMsg),
	}
}

func CreateCommitMsg(consensus cs.Consensus) msg.CommitMsg{
	return msg.CommitMsg{
		ConsensusID: consensus.ConsensusID,
	}
}
