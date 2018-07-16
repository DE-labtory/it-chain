package consensus

import (
	"github.com/it-chain/midgard"
)

// Publish part

type PrepareMsgAddedEvent struct {
	midgard.EventModel
	PrepareMsg struct {
		ConsensusId ConsensusId
		SenderId    string
		BlockHash   []byte
	}
}

type CommitMsgAddedEvent struct {
	midgard.EventModel
	CommitMsg struct {
		ConsensusId ConsensusId
		SenderId    string
	}
}

type ConsensusCreatedEvent struct {
	midgard.EventModel
	Consensus struct {
		ConsensusID     ConsensusId
		Representatives []*Representative
		Block           ProposedBlock
		CurrentState    State
		PrepareMsgPool  PrepareMsgPool
		CommitMsgPool   CommitMsgPool
	}
}

// Preprepare msg를 보냈을 때
type ConsensusPrePreparedEvent struct {
	midgard.EventModel
}

// Prepare msg를 보냈을 때
type ConsensusPreparedEvent struct {
	midgard.EventModel
}

// Commit msg를 보냈을 때
type ConsensusCommittedEvent struct {
	midgard.EventModel
}

// block 저장이 끝나 state가 idle이 될 때
type ConsensusFinishedEvent struct {
	midgard.EventModel
}

// Consume part

type LeaderChangedEvent struct {
	midgard.EventModel
	LeaderId string
}

type MemberJoinedEvent struct {
	midgard.EventModel
	MemberId string
}

type MemberRemovedEvent struct {
	midgard.EventModel
	MemberId string
}

// block이 저장되었을 때
type BlockSavedEvent struct {
	midgard.EventModel
}
