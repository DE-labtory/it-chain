package consensus

import (
	"github.com/it-chain/midgard"
)

// todo : Consensus로 시작하는 네이밍

// Publish part

type PrepareMsgAddedEvent struct {
	midgard.EventModel
	PrePrepareMsg struct {
		ConsensusId   ConsensusId
		SenderId      string
		ProposedBlock []byte
	}
}

type CommitMsgAddedEvent struct {
	midgard.EventModel
	CommitMsg struct {
		ConsensusId ConsensusId
		SenderId    string
	}
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

// todo : consensus를 위해 필요하지 않나? -> 고민해볼것

type ConsensusStartedEvent struct {
	midgard.EventModel
}

type PrepareFinishedEvent struct {
	midgard.EventModel
	ConsensusId string
}

type ConsensusFinishedEvent struct {
	midgard.EventModel
	ConsensusId string
}
