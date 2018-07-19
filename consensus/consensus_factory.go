package consensus

import (
	"github.com/rs/xid"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/consensus/infra/adapter"
)

// leader
func CreateConsensus(parliament Parliament, block ProposedBlock) (Consensus, error) {
	ps := adapter.NewParliamentService()
	representatives, err := ps.Elect(parliament)
	if err != nil {
		return Consensus{}, err
	}

	consensusID := NewConsensusId(xid.New().String())

	consensus := Consensus{
		ConsensusID:     consensusID,
		Representatives: representatives,
		Block:           block,
		CurrentState:    IDLE_STATE,
		PrepareMsgPool:  NewPrepareMsgPool(),
		CommitMsgPool:   NewCommitMsgPool(),
	}

	consensusCreatedEvent := ConsensusCreatedEvent{
		EventModel: midgard.EventModel{
			ID:      consensusID.Id,
		},
		Consensus: struct {
			ConsensusID     ConsensusId
			Representatives []*Representative
			Block           ProposedBlock
			CurrentState    State
			PrepareMsgPool  PrepareMsgPool
			CommitMsgPool   CommitMsgPool
		}{
			ConsensusID:     consensus.ConsensusID,
			Representatives: consensus.Representatives,
			Block:           consensus.Block,
			CurrentState:    consensus.CurrentState,
			PrepareMsgPool:  consensus.PrepareMsgPool,
			CommitMsgPool:   consensus.CommitMsgPool,
		},
	}

	if err := consensus.On(&consensusCreatedEvent); err != nil {
		return Consensus{}, err
	}

	if err := eventstore.Save(consensus.GetID(), consensusCreatedEvent); err != nil {
		return Consensus{}, err
	}

	return consensus, nil
}

// member
func ConstructConsensus(msg PrePrepareMsg) (Consensus, error) {
	consensus := Consensus{
		ConsensusID:     msg.ConsensusId,
		Representatives: msg.Representative,
		Block:           msg.ProposedBlock,
		CurrentState:    PREPREPARE_STATE,
		PrepareMsgPool:  NewPrepareMsgPool(),
		CommitMsgPool:   NewCommitMsgPool(),
	}

	consensusCreatedEvent := ConsensusCreatedEvent{
		EventModel: midgard.EventModel{
			ID:      consensus.GetID(),
		},
		Consensus: struct {
			ConsensusID     ConsensusId
			Representatives []*Representative
			Block           ProposedBlock
			CurrentState    State
			PrepareMsgPool  PrepareMsgPool
			CommitMsgPool   CommitMsgPool
		}{
			ConsensusID:     consensus.ConsensusID,
			Representatives: consensus.Representatives,
			Block:           consensus.Block,
			CurrentState:    consensus.CurrentState,
			PrepareMsgPool:  consensus.PrepareMsgPool,
			CommitMsgPool:   consensus.CommitMsgPool,
		},
	}

	if err := consensus.On(&consensusCreatedEvent); err != nil {
		return Consensus{}, err
	}

	if err := eventstore.Save(consensus.GetID(), consensusCreatedEvent); err != nil {
		return Consensus{}, err
	}

	return consensus, nil
}