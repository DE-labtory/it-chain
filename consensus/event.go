package consensus

import "github.com/it-chain/midgard"

// Publish part
type ConsensusMessagePublishedEvent struct {
	midgard.EventModel
	ConsensusMsg string
}

// todo : Blockchain 모듈 참고
type BlockCreatedEvent struct {
	midgard.EventModel
}
