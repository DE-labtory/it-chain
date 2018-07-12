package api

import (
	"github.com/it-chain/it-chain-Engine/consensus"
	"github.com/it-chain/midgard"
)

type ConsensusApi struct {
	eventRepository *midgard.Repository
}

func NewConsensusApi(eventRepository *midgard.Repository) ConsensusApi {
	return ConsensusApi{
		eventRepository: eventRepository,
	}
}

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
