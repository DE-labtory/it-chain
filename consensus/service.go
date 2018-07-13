package consensus

type Serializable interface {
	ToByte() ([]byte, error)
}

type MessageService interface {
	BroadcastMsg(Msg Serializable, representatives []*Representative)
	CreateConfirmedBlock(block ProposedBlock)
	IsLeaderMessage(msg PrePrepareMsg, leader Leader) bool
}

type ParliamentService interface {
}

type ElectionService interface {
}

type ConsensusService interface {
	ConstructConsensus(msg PrePrepareMsg)
}
