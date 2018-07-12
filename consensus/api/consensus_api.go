package api

import (
	"github.com/it-chain/it-chain-Engine/consensus"
	"github.com/it-chain/midgard"
)

type ConsensusApi struct {
	eventRepository *midgard.Repository
	consensus       consensus.Consensus
	parliament      consensus.Parliament
}

func NewConsensusApi(eventRepository *midgard.Repository, consensus consensus.Consensus, parliament consensus.Parliament) ConsensusApi {
	return ConsensusApi{
		eventRepository: eventRepository,
		consensus:       consensus,
		parliament:      parliament,
	}
}

func (cApi ConsensusApi) StartConsensus(userId consensus.MemberId, block consensus.ProposedBlock) error {
	parliament := cApi.parliament

	if parliament.IsNeedConsensus() {
		cApi.consensus.Start()

	} else {

	}

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
