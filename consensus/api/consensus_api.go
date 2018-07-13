package api

import (
	c "github.com/it-chain/it-chain-Engine/consensus"
	"github.com/it-chain/midgard"
)

type ConsensusApi struct {
	eventRepository *midgard.Repository
	consensus       c.Consensus
	parliament      c.Parliament
}

func NewConsensusApi(eventRepository *midgard.Repository, consensus c.Consensus, parliament c.Parliament) ConsensusApi {
	return ConsensusApi{
		eventRepository: eventRepository,
		consensus:       consensus,
		parliament:      parliament,
	}
}

// todo
func (cApi ConsensusApi) StartConsensus(userId c.MemberId, block c.ProposedBlock) error {
	//parliament := cApi.parliament
	//
	//if parliament.IsNeedConsensus() {
	//	consensus, err := c.CreateConsensus(parliament, block)
	//
	//	if err != nil {
	//		return err
	//	}
	//
	//	consensus.Start()
	//	cApi.consensus = *consensus
	//
	//	PrePrepareMsg := c.CreatePrePrepareMsg(*consensus)
	//
	//} else {
	//
	//}
	//
	return nil
}

func (cApi ConsensusApi) ReceivePrePrepareMsg(msg c.PrePrepareMsg) {
	return
}

func (cApi ConsensusApi) ReceivePrepareMsg(msg c.PrepareMsg) {
	return
}

func (cApi ConsensusApi) ReceiveCommitMsg(msg c.CommitMsg) {
	return
}
