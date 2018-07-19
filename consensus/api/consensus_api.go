package api

import (
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/midgard"
)

type ConsensusApi struct {
	eventRepository *midgard.Repository
}

// todo : Event Sourcing 첨가

func (cApi ConsensusApi) StartConsensus(userId consensus.MemberId, block consensus.ProposedBlock) error {
	return nil
}

func (cApi ConsensusApi) ReceivePrePrepareMsg(msg consensus.PrePrepareMsg) {
	return
}

func (cApi ConsensusApi) ReceivePrepareMsg(msg consensus.PrepareMsg) {
	return
}

func (cApi ConsensusApi) ReceiveCommitMsg(msg consensus.CommitMsg) {
	return
}
