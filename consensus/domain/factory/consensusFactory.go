package factory

import (
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/parliament"
	"github.com/it-chain/it-chain-Engine/consensus/domain/service"
	"github.com/rs/xid"
)

func CreateConsensus(parliament parliament.Parliament, block consensus.Block) (*consensus.Consensus, error) {

	//대표자정책에 의해 결정
	representatives, err := service.Elect(parliament)

	if err != nil {
		return nil, err
	}

	return &consensus.Consensus{
		ConsensusID:     consensus.NewConsensusID(xid.New().String()),
		Representatives: representatives,
		Block:           block,
		CurrentState:    new(consensus.IdleState),
	}, nil
}
