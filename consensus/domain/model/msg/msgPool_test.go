package msg

import (
	"testing"

	cs "github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
	"github.com/stretchr/testify/assert"
)

func TestMsgPool_InsertPrepareMsg(t *testing.T) {

	msgPool := &MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]CommitMsg),
	}

	prepareMsg := PrepareMsg{ConsensusID: cs.ConsensusID{"1"}}
	msgPool.InsertPrepareMsg(prepareMsg)

	assert.Equal(t, len(msgPool.PrepareMsgPool), 1)
	assert.Equal(t, msgPool.PrepareMsgPool[cs.ConsensusID{"1"}], []PrepareMsg{prepareMsg})
	assert.Equal(t, len(msgPool.PrepareMsgPool[cs.ConsensusID{"1"}]), 1)

	msgPool.InsertPrepareMsg(prepareMsg)

	assert.Equal(t, len(msgPool.PrepareMsgPool), 1)
	assert.Equal(t, 2, len(msgPool.PrepareMsgPool[cs.ConsensusID{"1"}]))
	assert.Equal(t, msgPool.PrepareMsgPool[cs.ConsensusID{"1"}], []PrepareMsg{prepareMsg, prepareMsg})
}

func TestMsgPool_FindPrepareMsgsByConsensusID(t *testing.T) {

	msgPool := &MsgPool{
		PrepareMsgPool: make(map[cs.ConsensusID][]PrepareMsg),
		CommitMsgPool:  make(map[cs.ConsensusID][]CommitMsg),
	}

	prepareMsg := PrepareMsg{ConsensusID: cs.ConsensusID{"1"}}
	msgPool.InsertPrepareMsg(prepareMsg)

	assert.Equal(t, 1, len(msgPool.FindPrepareMsgsByConsensusID(prepareMsg.ConsensusID)))

	prepareMsg2 := PrepareMsg{ConsensusID: cs.ConsensusID{"1"}}
	msgPool.InsertPrepareMsg(prepareMsg2)

	assert.Equal(t, 2, len(msgPool.FindPrepareMsgsByConsensusID(prepareMsg.ConsensusID)))
}
