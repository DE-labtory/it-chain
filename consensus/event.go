package consensus

import "github.com/it-chain/midgard"

// todo : Consensus로 시작하는 네이밍

// Publish part

type PrePrepareMsgCreatedEvent struct {
	midgard.EventModel
	PrePrepareMsg PrePrepareMsg
}

type PrepareMsgCreatedEvent struct {
	midgard.EventModel
	PrepareMsg PrepareMsg
}

type CommitMsgCreatedEvent struct {
	midgard.EventModel
	CommitMsg CommitMsg
}

// todo : Blockchain 모듈 참고
type BlockCreatedEvent struct {
	midgard.EventModel
}

// Consume part

type PrePrepareMsgArrivedEvent struct {
	midgard.EventModel
	PrePrepareMsg PrePrepareMsg
}

type PrepareMsgArrivedEvent struct {
	midgard.EventModel
	PrepareMsg PrepareMsg
}

type CommitMsgArrivedEvent struct {
	midgard.EventModel
	CommitMsg CommitMsg
}

type ConsensusStartedEvent struct {
	midgard.EventModel
}

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

type PrepareFinishedEvent struct {
	midgard.EventModel
	ConsensusId string
}

type ConsensusFinishedEvent struct {
	midgard.EventModel
	ConsensusId string
}
