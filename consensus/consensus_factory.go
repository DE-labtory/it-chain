package consensus

import (
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

// leader
func CreateConsensus(parliament Parliament, block ProposedBlock) (*Consensus, error) {
	representatives, err := Elect(parliament)
	if err != nil {
		return &Consensus{}, err
	}

	consensusID := NewConsensusId(xid.New().String())
	consensus := &Consensus{}
	consensusCreatedEvent := newConsensusCreatedEvent(consensusID, representatives, block)

	if err := OnAndSave(consensus, &consensusCreatedEvent); err != nil {
		return &Consensus{}, err
	}

	return consensus, nil
}

// member
func ConstructConsensus(msg PrePrepareMsg) (*Consensus, error) {
	consensus := &Consensus{}
	consensusCreatedEvent := newConsensusCreatedEvent(msg.ConsensusId, msg.Representative, msg.ProposedBlock)

	if err := OnAndSave(consensus, &consensusCreatedEvent); err != nil {
		return &Consensus{}, err
	}

	return consensus, nil
}

func newConsensusCreatedEvent(cID ConsensusId, r []*Representative, b ProposedBlock) ConsensusCreatedEvent {
	return ConsensusCreatedEvent{
		EventModel: midgard.EventModel{
			ID: cID.Id,
		},
		Consensus: struct {
			ConsensusID     ConsensusId
			Representatives []*Representative
			Block           ProposedBlock
			CurrentState    State
			PrepareMsgPool  PrepareMsgPool
			CommitMsgPool   CommitMsgPool
		}{
			ConsensusID:     cID,
			Representatives: r,
			Block:           b,
			CurrentState:    IDLE_STATE,
			PrepareMsgPool:  NewPrepareMsgPool(),
			CommitMsgPool:   NewCommitMsgPool(),
		},
	}
}
