package adapter

import c "github.com/it-chain/it-chain-Engine/consensus"

type ConsensusService struct {
	consensus *c.Consensus
}

func NewConsensusService(consensus c.Consensus) *ConsensusService {
	return &ConsensusService{
		consensus: &consensus,
	}
}

func (c ConsensusService) ConstructConsensus(msg c.PrePrepareMsg) {
	c.consensus.ConsensusID = msg.ConsensusId
	c.consensus.Representatives = msg.Representative
	c.consensus.Block = msg.ProposedBlock
}
