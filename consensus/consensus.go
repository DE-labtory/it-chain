package consensus

import (
	"errors"
	"fmt"

	"github.com/it-chain/midgard"
)

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
	Representatives []*Representative
	Block           ProposedBlock
	CurrentState    State
	PrepareMsgPool  PrepareMsgPool
	CommitMsgPool   CommitMsgPool
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

func (c *Consensus) ToIdleState() {
	c.CurrentState = IDLE_STATE
}

func (c *Consensus) SavePrepareMsg(prepareMsg *PrepareMsg) (*PrepareMsgAddedEvent, error) {
	if c.ConsensusID.Id != prepareMsg.ConsensusId.Id {
		return nil, errors.New("Consensus ID is not same")
	}

	prepareMsgAddedEvent := PrepareMsgAddedEvent{
		EventModel: midgard.EventModel{
			ID: c.ConsensusID.Id,
		},
		PrepareMsg: struct {
			ConsensusId   ConsensusId
			SenderId      string
			ProposedBlock []byte
		}{ConsensusId: prepareMsg.ConsensusId, SenderId: prepareMsg.SenderId, ProposedBlock: prepareMsg.ProposedBlock},
	}

	c.On(&prepareMsgAddedEvent)

	return &prepareMsgAddedEvent, nil
}

func (c *Consensus) SaveCommitMsg(commitMsg CommitMsg) (*CommitMsgAddedEvent, error) {
	if c.ConsensusID.Id != commitMsg.ConsensusId.Id {
		return nil, errors.New("Consensus ID is not same")
	}

	commitMsgAddedEvent := CommitMsgAddedEvent{
		EventModel: midgard.EventModel{
			ID: c.ConsensusID.Id,
		},
		CommitMsg: struct {
			ConsensusId ConsensusId
			SenderId    string
		}{ConsensusId: commitMsg.ConsensusId, SenderId: commitMsg.SenderId},
	}

	c.On(&commitMsgAddedEvent)

	return &commitMsgAddedEvent, nil
}

func (c *Consensus) On(event midgard.Event) error {
	switch v := event.(type) {

	case *PrepareMsgAddedEvent:
		c.PrepareMsgPool.Save(&PrepareMsg{
			ConsensusId:   v.PrepareMsg.ConsensusId,
			SenderId:      v.PrepareMsg.SenderId,
			ProposedBlock: v.PrepareMsg.ProposedBlock,
		})

	case *CommitMsgAddedEvent:
		c.CommitMsgPool.Save(&CommitMsg{
			ConsensusId: v.CommitMsg.ConsensusId,
			SenderId:    v.CommitMsg.SenderId,
		})

	case *ConsensusStartedEvent:
		c.Start()

	case *ConsensusPreparedEvent:
		c.ToCommitState()

	case *ConsensusFinishedEvent:
		c.ToIdleState()

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}
