package consensus

type State string

const (
	IDLE_STATE       State = "IdleState"
	PREPREPARE_STATE State = "PrePrepareState"
	PREPARE_STATE    State = "PrepareState"
	COMMIT_STATE     State = "CommitState"
)

type ConsensusId struct {
	Id string
}

func NewConsensusId(id string) ConsensusId {
	return ConsensusId{
		Id: id,
	}
}

type Consensus struct {
	ConsensusID     ConsensusId
	// todo : what is "Representativ"?
	//Representatives []*Representative
	Block           ProposedBlock
	CurrentState    State
}

func (c *Consensus) Start() {
	c.CurrentState = PREPARE_STATE
}

func (c *Consensus) IsPrepareState() bool {

	if c.CurrentState == PREPARE_STATE {
		return true
	}
	return false
}

func (c *Consensus) IsCommitState() bool {

	if c.CurrentState == COMMIT_STATE {
		return true
	}
	return false
}

func (c *Consensus) ToCommitState() {
	c.CurrentState = COMMIT_STATE
}
