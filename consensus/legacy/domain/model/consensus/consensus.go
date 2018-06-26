package consensus

type State string

const (
	IDLE_STATE       State = "IdleState"
	PREPREPARE_STATE State = "PreprepareState"
	PREPARE_STATE    State = "PrepareState"
	COMMIT_STATE     State = "CommitState"
)

type ConsensusID struct {
	ID string
}

func NewConsensusID(id string) ConsensusID {
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
