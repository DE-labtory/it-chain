package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareMsgPool_Save(t *testing.T) {
	// given
	pPool := NewPrepareMsgPool()

	// case 1 : save
	pMsg := PrepareMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s1",
		BlockHash:   make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 1, len(pPool.messages))

	// case 2 : save
	pMsg = PrepareMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s2",
		BlockHash:   make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.messages))

	// case 3 : same sender
	pMsg = PrepareMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s2",
		BlockHash:   make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.messages))

	// case 4 : block hash is is nil
	pMsg = PrepareMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s3",
		BlockHash:   nil,
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.messages))
}

func TestPrepareMsgPool_RemoveAllMsgs(t *testing.T) {
	// given
	pPool := NewPrepareMsgPool()

	pMsg := PrepareMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s1",
		BlockHash:   make([]byte, 0),
	}

	pPool.Save(&pMsg)

	// when
	pPool.RemoveAllMsgs()

	// then
	assert.Equal(t, 0, len(pPool.messages))
}

func TestCommitMsgPool_Save(t *testing.T) {
	// given
	cPool := NewCommitMsgPool()

	// case 1 : save
	cMsg := CommitMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s1",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 1, len(cPool.messages))

	// case 2 : save
	cMsg = CommitMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s2",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 2, len(cPool.messages))

	// case 3 : same sender
	cMsg = CommitMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s2",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 2, len(cPool.messages))
}

func TestCommitMsgPool_RemoveAllMsgs(t *testing.T) {
	// given
	cPool := NewCommitMsgPool()

	cMsg := CommitMsg{
		ConsensusId: ConsensusId{"c1"},
		SenderId:    "s1",
	}

	cPool.Save(&cMsg)

	// when
	cPool.RemoveAllMsgs()

	// then
	assert.Equal(t, 0, len(cPool.messages))
}

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
