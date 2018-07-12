package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsensus_SavePrepareMsg(t *testing.T) {
	// given
	c := Consensus{
		ConsensusID:     NewConsensusId("c1"),
		Representatives: nil,
		Block: ProposedBlock{
			Seal: make([]byte, 0),
		},
		CurrentState:   IDLE_STATE,
		PrepareMsgPool: NewPrepareMsgPool(),
		CommitMsgPool:  NewCommitMsgPool(),
	}

	// case 1 : save
	pMsg := PrepareMsg{
		ConsensusId: NewConsensusId("c1"),
		SenderId:    "s1",
		BlockHash:   make([]byte, 0),
	}

	// when
	event, _ := c.SavePrepareMsg(&pMsg)

	// then
	assert.Equal(t, 1, len(c.PrepareMsgPool.messages))
	assert.Equal(t, c.ConsensusID.Id, event.ID)

	// case 2 : incorrect consensus ID
	pMsg = PrepareMsg{
		ConsensusId: NewConsensusId("c2"),
		SenderId:    "s1",
		BlockHash:   make([]byte, 0),
	}

	// when
	c.SavePrepareMsg(&pMsg)

	//then
	assert.Equal(t, 1, len(c.PrepareMsgPool.messages))
}

func TestConsensus_SaveCommitMsg(t *testing.T) {
	// given
	c := Consensus{
		ConsensusID:     NewConsensusId("c1"),
		Representatives: nil,
		Block: ProposedBlock{
			Seal: make([]byte, 0),
		},
		CurrentState:   IDLE_STATE,
		PrepareMsgPool: NewPrepareMsgPool(),
		CommitMsgPool:  NewCommitMsgPool(),
	}

	// case 1 : save
	cMsg := CommitMsg{
		ConsensusId: NewConsensusId("c1"),
		SenderId:    "s1",
	}

	// when
	event, _ := c.SaveCommitMsg(&cMsg)

	// then
	assert.Equal(t, 1, len(c.CommitMsgPool.messages))
	assert.Equal(t, c.ConsensusID.Id, event.ID)

	// case 2 : incorrect consensus ID
	cMsg = CommitMsg{
		ConsensusId: NewConsensusId("c2"),
		SenderId:    "s1",
	}

	// when
	c.SaveCommitMsg(&cMsg)

	//then
	assert.Equal(t, 1, len(c.CommitMsgPool.messages))
}
