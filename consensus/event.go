package consensus

import "github.com/it-chain/midgard"

// todo : Consensus로 시작하는 네이밍

// Publish part

type ConsensusMessagePublishedEvent struct {
	midgard.EventModel
	ConsensusMsg string
}

// todo : Blockchain 모듈 참고
type BlockCreatedEvent struct {
	midgard.EventModel
}

// Consume part

type ConsensusMessageArrivedEvent struct {
	midgard.EventModel
	ConsensusMsg string
}

type ConsensusStartedEvent struct {
	midgard.EventModel
}

type LeaderChangedEvent struct {
	midgard.EventModel
}