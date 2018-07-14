package api

import (
	c "github.com/it-chain/it-chain-Engine/consensus"
	"github.com/it-chain/midgard"
)

type ConsensusApi struct {
	eventRepository  *midgard.Repository
	msgService       c.MessageService
	consensusService c.ConsensusService
}

func NewConsensusApi(eventRepository *midgard.Repository, consensus *c.Consensus, parliament *c.Parliament) ConsensusApi {
	return ConsensusApi{
		eventRepository: eventRepository,
	}
}

func (cApi ConsensusApi) StartConsensus(userId c.MemberId, block c.ProposedBlock) error {
	parliament := c.NewParliament()
	cApi.eventRepository.Load(parliament, c.PARLIAMENT_AID)

	if parliament.IsNeedConsensus() {
		consensus, err := c.CreateConsensus(*parliament, block)

		if err != nil {
			return err
		}

		consensus.Start()
		cApi.consensus = consensus

		PrePrepareMsg := c.CreatePrePrepareMsg(*consensus)
		cApi.msgService.BroadcastMsg(PrePrepareMsg, consensus.Representatives)
	} else {
		// config 따라 다름...
	}

	return nil
}

func (cApi ConsensusApi) ReceivePrePrepareMsg(msg c.PrePrepareMsg) {
	mService := cApi.msgService
	parliament := cApi.parliament
	consensus := cApi.consensus

	if mService.IsLeaderMessage(msg, *parliament.Leader) {
		cService := cApi.consensusService

		cService.ConstructConsensus(msg)
		cApi.consensus.Start()

		PrepareMsg := c.CreatePrepareMsg(*consensus)
		cApi.msgService.BroadcastMsg(PrepareMsg, consensus.Representatives)
		consensus.Start()
	} else {
		return
	}
}

func (cApi ConsensusApi) ReceivePrepareMsg(msg c.PrepareMsg) {
	cApi.consensus.SavePrepareMsg(&msg)
	consensus := cApi.consensus

	if c.CheckConsensusPolicy(*consensus) {
		CommitMsg := c.CreateCommitMsg(*consensus)
		cApi.msgService.BroadcastMsg(CommitMsg, consensus.Representatives)
		consensus.ToPrepareState()
	} else {
		return
	}
}

func (cApi ConsensusApi) ReceiveCommitMsg(msg c.CommitMsg) {
	return
}
