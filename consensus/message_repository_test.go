package consensus

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPrepareMsgRepositoryImpl_InsertPrepareMsg(t *testing.T) {
	// given
	prepareMsgRepository := NewPrepareMsgRepository()

	prepareMsg := PrepareMsg{
		ConsensusId: ConsensusId{"1"},
		SenderId:    "1",
	}

	// when
	prepareMsgRepository.InsertPrepareMsg(prepareMsg)

	// then
	assert.Equal(t, 1, len(prepareMsgRepository.PrepareMsgPool))
	assert.Equal(t, []PrepareMsg{prepareMsg}, prepareMsgRepository.PrepareMsgPool[ConsensusId{"1"}])

	// when : Same representative sends same message
	prepareMsgRepository.InsertPrepareMsg(prepareMsg)

	// then
	assert.Equal(t, 1, len(prepareMsgRepository.PrepareMsgPool))
	assert.Equal(t, []PrepareMsg{prepareMsg}, prepareMsgRepository.PrepareMsgPool[ConsensusId{"1"}])

	// given
	prepareMsg2 := PrepareMsg{
		ConsensusId: ConsensusId{"2"},
		SenderId:    "1",
	}

	// when
	prepareMsgRepository.InsertPrepareMsg(prepareMsg2)

	// then
	assert.Equal(t, 2, len(prepareMsgRepository.PrepareMsgPool))
	assert.Equal(t, []PrepareMsg{prepareMsg}, prepareMsgRepository.PrepareMsgPool[ConsensusId{"1"}])
	assert.Equal(t, []PrepareMsg{prepareMsg2}, prepareMsgRepository.PrepareMsgPool[ConsensusId{"2"}])
}

func TestCommitMsgRepositoryImpl_InsertCommitMsg(t *testing.T) {
	// given
	commitMsgRepository := NewCommitMsgRepository()

	commitMsg := CommitMsg{
		ConsensusId: ConsensusId{"1"},
		SenderId:    "1",
	}

	// when
	commitMsgRepository.InsertCommitMsg(commitMsg)

	// then
	assert.Equal(t, 1, len(commitMsgRepository.CommitMsgPool))
	assert.Equal(t, []CommitMsg{commitMsg}, commitMsgRepository.CommitMsgPool[ConsensusId{"1"}])

	// when : Same representative sends same message
	commitMsgRepository.InsertCommitMsg(commitMsg)

	// then
	assert.Equal(t, 1, len(commitMsgRepository.CommitMsgPool))
	assert.Equal(t, []CommitMsg{commitMsg}, commitMsgRepository.CommitMsgPool[ConsensusId{"1"}])

	// given
	commitMsg2 := CommitMsg{
		ConsensusId: ConsensusId{"2"},
		SenderId:    "1",
	}

	// when
	commitMsgRepository.InsertCommitMsg(commitMsg2)

	// then
	assert.Equal(t, 2, len(commitMsgRepository.CommitMsgPool))
	assert.Equal(t, []CommitMsg{commitMsg}, commitMsgRepository.CommitMsgPool[ConsensusId{"1"}])
	assert.Equal(t, []CommitMsg{commitMsg2}, commitMsgRepository.CommitMsgPool[ConsensusId{"2"}])
}

func TestPrepareMsgRepositoryImpl_DeleteAllPrepareMsg(t *testing.T) {
	// given
	prepareMsgRepository := NewPrepareMsgRepository()

	prepareMsg := PrepareMsg{
		ConsensusId: ConsensusId{"1"},
		SenderId:    "1",
	}

	prepareMsgRepository.InsertPrepareMsg(prepareMsg)

	// when
	prepareMsgRepository.DeleteAllPrepareMsg(prepareMsg.ConsensusId)

	// then
	assert.Equal(t, 0, len(prepareMsgRepository.PrepareMsgPool))
	assert.Nil(t, prepareMsgRepository.PrepareMsgPool[prepareMsg.ConsensusId])
}

func TestCommitMsgRepositoryImpl_DeleteAllCommitMsg(t *testing.T) {
	// given
	commitMsgRepository := NewCommitMsgRepository()

	commitMsg := CommitMsg{
		ConsensusId: ConsensusId{"1"},
		SenderId:    "1",
	}

	commitMsgRepository.InsertCommitMsg(commitMsg)

	// when
	commitMsgRepository.DeleteAllCommitMsg(commitMsg.ConsensusId)

	// then
	assert.Equal(t, 0, len(commitMsgRepository.CommitMsgPool))
	assert.Nil(t, commitMsgRepository.CommitMsgPool[commitMsg.ConsensusId])
}

func TestPrepareMsgRepositoryImpl_FindPrepareMsgsByConsensusID(t *testing.T) {
	// given
	prepareMsgRepository := NewPrepareMsgRepository()

	prepareMsg := PrepareMsg{
		ConsensusId: ConsensusId{"1"},
		SenderId:    "1",
	}

	prepareMsg2 := PrepareMsg{
		ConsensusId: ConsensusId{"1"},
		SenderId:    "2",
	}

	prepareMsg3 := PrepareMsg{
		ConsensusId: ConsensusId{"2"},
		SenderId:    "1",
	}

	prepareMsgRepository.InsertPrepareMsg(prepareMsg)
	prepareMsgRepository.InsertPrepareMsg(prepareMsg2)
	prepareMsgRepository.InsertPrepareMsg(prepareMsg3)

	// when
	msgs := prepareMsgRepository.FindPrepareMsgsByConsensusID(prepareMsg.ConsensusId)

	// then
	assert.Equal(t, 2, len(msgs))

	// when
	msgs = prepareMsgRepository.FindPrepareMsgsByConsensusID(prepareMsg3.ConsensusId)

	// then
	assert.Equal(t, 1, len(msgs))
}

func TestCommitMsgRepositoryImpl_FindCommitMsgsByConsensusID(t *testing.T) {
	// given
	commitMsgRepository := NewCommitMsgRepository()

	commitMsg := CommitMsg{
		ConsensusId: ConsensusId{"1"},
		SenderId:    "1",
	}

	commitMsg2 := CommitMsg{
		ConsensusId: ConsensusId{"1"},
		SenderId:    "2",
	}

	commitMsg3 := CommitMsg{
		ConsensusId: ConsensusId{"2"},
		SenderId:    "1",
	}

	commitMsgRepository.InsertCommitMsg(commitMsg)
	commitMsgRepository.InsertCommitMsg(commitMsg2)
	commitMsgRepository.InsertCommitMsg(commitMsg3)

	// when
	msgs := commitMsgRepository.FindCommitMsgsByConsensusID(commitMsg.ConsensusId)

	// then
	assert.Equal(t, 2, len(msgs))

	// when
	msgs = commitMsgRepository.FindCommitMsgsByConsensusID(commitMsg3.ConsensusId)

	// then
	assert.Equal(t, 1, len(msgs))
}
