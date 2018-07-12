package consensus

type Serializable interface {
	ToByte() ([]byte, error)
}

type MessageService interface {
	BroadcastMsg(Msg Serializable, representatives []*Representative)
}

type ParliamentService interface {
}

type ElectionService interface {
}

type ConsensusService interface {
}
