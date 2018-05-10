package event

//type Sendable struct {
//	Ids  []string
//	Data []byte
//}

type ConsensusMessageType int

const (
	PREPREPARE ConsensusMessageType = 0
	PREPARE    ConsensusMessageType = 1
	COMMIT     ConsensusMessageType = 2
)

type ConsensusMessagePublishEvent struct {
	Ids  []string
	Data []byte
}

//todo define event
type BlockConfirmEvent struct {
	Block []byte
}

type ConsensusCreateCmd struct {
	UserID string
	Block  []byte
}
