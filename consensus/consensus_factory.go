package consensus

import "github.com/rs/xid"

func CreateConsensus(parliament Parliament, block ProposedBlock) (*Consensus, error) {
	representatives, err := Elect(parliament)
	if err != nil {
		return nil, err
	}

	return &Consensus{
		ConsensusID:     ConsensusId{xid.New().String()},
		Representatives: representatives,
		Block:           block,
		CurrentState:    IDLE_STATE,
		PrepareMsgPool:  NewPrepareMsgPool(),
		CommitMsgPool:   NewCommitMsgPool(),
	}, nil
}
