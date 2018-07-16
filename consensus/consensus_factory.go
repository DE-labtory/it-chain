package consensus

import (
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

func CreateConsensus(parliament Parliament, block ProposedBlock) (*Consensus, error) {
	representatives, err := Elect(parliament)
	if err != nil {
		return nil, err
	}

	return &Consensus{
		AggregateModel: midgard.AggregateModel{
			ID: CONSENSUS_AID,
		},
		ConsensusID:     ConsensusId{xid.New().String()},
		Representatives: representatives,
		Block:           block,
		CurrentState:    IDLE_STATE,
		PrepareMsgPool:  NewPrepareMsgPool(),
		CommitMsgPool:   NewCommitMsgPool(),
	}, nil
}

func CreatePrePrepareMsg(consensus Consensus) PrePrepareMsg {
	return PrePrepareMsg{
		ConsensusId:    consensus.ConsensusID,
		Representative: consensus.Representatives,
		ProposedBlock:  consensus.Block,
	}
}

func CreatePrepareMsg(consensus Consensus) PrepareMsg {
	return PrepareMsg{
		ConsensusId: consensus.ConsensusID,
		BlockHash:   consensus.Block.Seal,
	}
}

func CreateCommitMsg(consensus Consensus) CommitMsg {
	return CommitMsg{
		ConsensusId: consensus.ConsensusID,
	}
}
