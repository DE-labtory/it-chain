package event_message

type Sendable struct {
	Ids  []string
	Data []byte
}

type ConsensusMessageType int

const (
	PREPREPARE ConsensusMessageType = 0
	PREPARE    ConsensusMessageType = 1
	COMMIT     ConsensusMessageType = 2
)

type StartConsensusEvent struct {
	Block  []byte
	UserID string
}

type ReceviedConsensusMessageEvent struct {
	MessageType ConsensusMessageType
	MessageBody []byte
}
