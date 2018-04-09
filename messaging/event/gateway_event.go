package event

type ConsensusStartEvent struct {
	Block  []byte
	UserID string
}

type ConsensusMessageArriveEvent struct {
	MessageType ConsensusMessageType
	MessageBody []byte
}
