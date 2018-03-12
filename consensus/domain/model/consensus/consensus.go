package consensus

type ConsensusID struct{
	ID string
}

func NewConsensusID (id string) ConsensusID{
	return ConsensusID{
		ID: id,
	}
}

type LeaderID string

type Consensus struct {
	ConsensusID     ConsensusID
	Representatives []*Representative
	Block           Block
	CurrentState    State
}

func (c *Consensus) Start(){
	c.CurrentState = new(PrepareState)
}