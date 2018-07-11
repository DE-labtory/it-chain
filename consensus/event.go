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

// Preprepare msg를 받았을 때
type ConsensusStartedEvent struct {
	midgard.EventModel
}

// Prepare msg를 받아서 commit msg를 받는 상태가 될 때
type ConsensusPreparedEvent struct {
	midgard.EventModel
	ConsensusId string
}

// Commit msg를 받아서 consensus가 끝났을 때
type ConsensusFinishedEvent struct {
	midgard.EventModel
	ConsensusId string
}
