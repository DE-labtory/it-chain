package api

import (
	c "github.com/it-chain/it-chain-Engine/consensus"
	"github.com/it-chain/midgard"
)

type ConsensusApi struct {
	eventRepository  *midgard.Repository
	consensus        c.Consensus
	parliament       c.Parliament
	msgService       c.MessageService
	consensusService c.ConsensusService
}

func NewConsensusApi(eventRepository *midgard.Repository, consensus c.Consensus, parliament c.Parliament) ConsensusApi {
	return ConsensusApi{
		eventRepository: eventRepository,
		consensus:       consensus,
		parliament:      parliament,
	}
}

func (cApi ConsensusApi) StartConsensus(userId c.MemberId, block c.ProposedBlock) error {
	parliament := cApi.parliament

	if parliament.IsNeedConsensus() {
		consensus, err := c.CreateConsensus(parliament, block)

		if err != nil {
			return err
		}

		consensus.Start()
		cApi.consensus = *consensus

		PrePrepareMsg := c.CreatePrePrepareMsg(*consensus)
		// BroadcastMsg(PrePrepareMsg)
	} else {
		// config 따라 다름...
	}

	return nil
}

func (cApi ConsensusApi) ReceivePrePrepareMsg(msg c.PrePrepareMsg) {
	mService := cApi.msgService
	parliament := cApi.parliament

	if mService.IsLeaderMessage(msg, *parliament.Leader) {
		cService := cApi.consensusService

		cService.ConstructConsensus(msg)
		cApi.consensus.Start()

		PrepareMsg := c.CreatePrepareMsg(cApi.consensus)
		// BroadcastMsg(PrepareMsg)
		cApi.consensus.ToPrepareState()
	} else {
		return
	}
}

func (cApi ConsensusApi) ReceivePrepareMsg(msg c.PrepareMsg) {
	return
}

func (cApi ConsensusApi) ReceiveCommitMsg(msg c.CommitMsg) {
	return
}
