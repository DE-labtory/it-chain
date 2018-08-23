package adapter

import "github.com/it-chain/engine/blockchain"

type ConsensusService struct{}

func NewConsensusService() *ConsensusService {
	return &ConsensusService{}
}

// TODO
func (s *ConsensusService) ConsentBlock(block blockchain.DefaultBlock) error {
	return nil
}
