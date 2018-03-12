package repository

import "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"

type ConsensusRepository interface{
	Save(consensus consensus.Consensus)
	FindById(consensusId consensus.ConsensusID) *consensus.Consensus
}
