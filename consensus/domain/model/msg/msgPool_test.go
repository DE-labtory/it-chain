package msg

import (
	"testing"

	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/stretchr/testify/assert"
)

func TestMsgPool_InsertPrepareMsg(t *testing.T) {

	//given
	msgPool := &MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]CommitMsg),
	}

	prepareMsg := PrepareMsg{ConsensusID: cs.ConsensusID{"1"}, SenderID: "1"}
	//when
	msgPool.InsertPrepareMsg(prepareMsg)

	//then
	assert.Equal(t, len(msgPool.PrepareMsgPool), 1)
	assert.Equal(t, msgPool.PrepareMsgPool[cs.ConsensusID{"1"}], []PrepareMsg{prepareMsg})
	assert.Equal(t, len(msgPool.PrepareMsgPool[cs.ConsensusID{"1"}]), 1)

	//when
	msgPool.InsertPrepareMsg(prepareMsg)

	//then
	assert.Equal(t, len(msgPool.PrepareMsgPool), 1)
	assert.Equal(t, 1, len(msgPool.PrepareMsgPool[cs.ConsensusID{"1"}]))
	assert.Equal(t, msgPool.PrepareMsgPool[cs.ConsensusID{"1"}], []PrepareMsg{prepareMsg})

	//when
	prepareMsg2 := PrepareMsg{ConsensusID: cs.ConsensusID{"1"}, SenderID: "2"}
	msgPool.InsertPrepareMsg(prepareMsg2)

	//then
	assert.Equal(t, len(msgPool.PrepareMsgPool), 1)
	assert.Equal(t, 2, len(msgPool.PrepareMsgPool[cs.ConsensusID{"1"}]))
	assert.Equal(t, msgPool.PrepareMsgPool[cs.ConsensusID{"1"}], []PrepareMsg{prepareMsg, prepareMsg2})
}

func TestMsgPool_FindPrepareMsgsByConsensusID(t *testing.T) {

	//given
	msgPool := &MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]CommitMsg),
	}
	prepareMsg := PrepareMsg{ConsensusID: cs.ConsensusID{"1"}, SenderID: "2"}
	prepareMsg2 := PrepareMsg{ConsensusID: cs.ConsensusID{"1"}, SenderID: "3"}

	//when
	msgPool.InsertPrepareMsg(prepareMsg)

	//then
	assert.Equal(t, 1, len(msgPool.FindPrepareMsgsByConsensusID(prepareMsg.ConsensusID)))

	//when
	msgPool.InsertPrepareMsg(prepareMsg2)

	//then
	assert.Equal(t, 2, len(msgPool.FindPrepareMsgsByConsensusID(prepareMsg.ConsensusID)))
}

//commitMsg Insert Test
func TestMsgPool_InsertCommitMsg(t *testing.T) {

	//given
	msgPool := &MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]CommitMsg),
	}

	commitMsg := CommitMsg{ConsensusID: cs.ConsensusID{"1"}, SenderID: "1"}

	//when
	msgPool.InsertCommitMsg(commitMsg)

	//then
	assert.Equal(t, len(msgPool.CommitMsgPool), 1)
	assert.Equal(t, msgPool.CommitMsgPool[cs.ConsensusID{"1"}], []CommitMsg{commitMsg})
	assert.Equal(t, len(msgPool.CommitMsgPool[cs.ConsensusID{"1"}]), 1)

	//when 중복 sender 테스트
	msgPool.InsertCommitMsg(commitMsg)

	//then
	assert.Equal(t, len(msgPool.CommitMsgPool), 1)
	assert.Equal(t, 1, len(msgPool.CommitMsgPool[cs.ConsensusID{"1"}]))
	assert.Equal(t, msgPool.CommitMsgPool[cs.ConsensusID{"1"}], []CommitMsg{commitMsg})

	//when
	commitMsg2 := CommitMsg{ConsensusID: cs.ConsensusID{"1"}, SenderID: "2"}
	msgPool.InsertCommitMsg(commitMsg2)

	//then
	assert.Equal(t, len(msgPool.CommitMsgPool), 1)
	assert.Equal(t, 2, len(msgPool.CommitMsgPool[cs.ConsensusID{"1"}]))
	assert.Equal(t, msgPool.CommitMsgPool[cs.ConsensusID{"1"}], []CommitMsg{commitMsg, commitMsg2})
}

func TestMsgPool_DeleteCommitMsg(t *testing.T) {
	//given
	msgPool := &MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]CommitMsg),
	}

	commitMsg := CommitMsg{ConsensusID: cs.ConsensusID{"1"}}
	msgPool.InsertCommitMsg(commitMsg)
	msgPool.InsertCommitMsg(commitMsg)

	//when
	msgPool.DeleteCommitMsg(commitMsg.ConsensusID)

	//then
	assert.Equal(t, 0, len(msgPool.CommitMsgPool))
	assert.Nil(t, msgPool.CommitMsgPool[commitMsg.ConsensusID])
}

func TestMsgPool_DeletePrepareMsg(t *testing.T) {

	//given
	msgPool := &MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]CommitMsg),
	}

	prepareMsg := PrepareMsg{ConsensusID: cs.ConsensusID{"1"}}
	msgPool.InsertPrepareMsg(prepareMsg)
	msgPool.InsertPrepareMsg(prepareMsg)

	//when
	msgPool.DeletePrepareMsg(prepareMsg.ConsensusID)

	//then
	assert.Equal(t, 0, len(msgPool.CommitMsgPool))
	assert.Nil(t, msgPool.PrepareMsgPool[prepareMsg.ConsensusID])
}

func TestMsgPool_FindCommitMsgsByConsensusID(t *testing.T) {

	//given
	msgPool := &MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]CommitMsg),
	}

	commitMsg := CommitMsg{ConsensusID: cs.ConsensusID{"1"}, SenderID: "1"}

	//when
	msgPool.InsertCommitMsg(commitMsg)

	//then
	assert.Equal(t, 1, len(msgPool.FindCommitMsgsByConsensusID(commitMsg.ConsensusID)))
}
