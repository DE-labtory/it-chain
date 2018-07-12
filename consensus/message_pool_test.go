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
		ConsensusId:   ConsensusId{"c1"},
		SenderId:      "s1",
		ProposedBlock: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 1, len(pPool.messages))

	// case 2 : save
	pMsg = PrepareMsg{
		ConsensusId:   ConsensusId{"c1"},
		SenderId:      "s2",
		ProposedBlock: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.messages))

	// case 3 : same sender
	pMsg = PrepareMsg{
		ConsensusId:   ConsensusId{"c1"},
		SenderId:      "s2",
		ProposedBlock: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.messages))

	// case 4 : proposed block is is nil
	pMsg = PrepareMsg{
		ConsensusId:   ConsensusId{"c1"},
		SenderId:      "s3",
		ProposedBlock: nil,
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
		ConsensusId:   ConsensusId{"c1"},
		SenderId:      "s1",
		ProposedBlock: make([]byte, 0),
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
