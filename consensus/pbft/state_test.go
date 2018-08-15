/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pbft_test

import (
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/stretchr/testify/assert"
)

func TestPrepareMsgPool_Save(t *testing.T) {
	// given
	pPool := pbft.NewPrepareMsgPool()

	// case 1 : save
	pMsg := pbft.PrepareMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 1, len(pPool.Get()))

	// case 2 : save
	pMsg = pbft.PrepareMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s2",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.Get()))

	// case 3 : same sender
	pMsg = pbft.PrepareMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s2",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.Get()))

	// case 4 : block hash is is nil
	pMsg = pbft.PrepareMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s3",
		BlockHash: nil,
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.Get()))
}

func TestPrepareMsgPool_RemoveAllMsgs(t *testing.T) {
	// given
	pPool := pbft.NewPrepareMsgPool()

	pMsg := pbft.PrepareMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	pPool.Save(&pMsg)

	// when
	pPool.RemoveAllMsgs()

	// then
	assert.Equal(t, 0, len(pPool.Get()))
}

func TestCommitMsgPool_Save(t *testing.T) {
	// given
	cPool := pbft.NewCommitMsgPool()

	// case 1 : save
	cMsg := pbft.CommitMsg{
		StateID:  pbft.StateID{"c1"},
		SenderID: "s1",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 1, len(cPool.Get()))

	// case 2 : save
	cMsg = pbft.CommitMsg{
		StateID:  pbft.StateID{"c1"},
		SenderID: "s2",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 2, len(cPool.Get()))

	// case 3 : same sender
	cMsg = pbft.CommitMsg{
		StateID:  pbft.StateID{"c1"},
		SenderID: "s2",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 2, len(cPool.Get()))
}

func TestCommitMsgPool_RemoveAllMsgs(t *testing.T) {
	// given
	cPool := pbft.NewCommitMsgPool()

	cMsg := pbft.CommitMsg{
		StateID:  pbft.StateID{"c1"},
		SenderID: "s1",
	}

	cPool.Save(&cMsg)

	// when
	cPool.RemoveAllMsgs()

	// then
	assert.Equal(t, 0, len(cPool.Get()))
}

func TestConsensus_SavePrepareMsg(t *testing.T) {
	// given
	c := pbft.State{
		StateID:         pbft.NewStateID("c1"),
		Representatives: nil,
		Block: pbft.ProposedBlock{
			Seal: make([]byte, 0),
		},
		CurrentStage:   pbft.IDLE_STAGE,
		PrepareMsgPool: pbft.NewPrepareMsgPool(),
		CommitMsgPool:  pbft.NewCommitMsgPool(),
	}

	// case 1 : save
	pMsg := &pbft.PrepareMsg{
		StateID:   pbft.NewStateID("c1"),
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	err := c.SavePrepareMsg(pMsg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.PrepareMsgPool.Get()))
	assert.Equal(t, *pMsg, c.PrepareMsgPool.Get()[0])

	// case 2 : incorrect consensus ID
	pMsg = &pbft.PrepareMsg{
		StateID:   pbft.NewStateID("c2"),
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	err = c.SavePrepareMsg(pMsg)

	//then
	assert.Error(t, err)
	assert.Equal(t, 1, len(c.PrepareMsgPool.Get()))
}

func TestConsensus_SaveCommitMsg(t *testing.T) {
	// given
	c := pbft.State{
		StateID:         pbft.NewStateID("c1"),
		Representatives: nil,
		Block: pbft.ProposedBlock{
			Seal: make([]byte, 0),
		},
		CurrentStage:   pbft.IDLE_STAGE,
		PrepareMsgPool: pbft.NewPrepareMsgPool(),
		CommitMsgPool:  pbft.NewCommitMsgPool(),
	}

	// case 1 : save
	cMsg := &pbft.CommitMsg{
		StateID:  pbft.NewStateID("c1"),
		SenderID: "s1",
	}

	// when
	err := c.SaveCommitMsg(cMsg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.CommitMsgPool.Get()))
	assert.Equal(t, *cMsg, c.CommitMsgPool.Get()[0])

	// case 2 : incorrect consensus ID
	cMsg = &pbft.CommitMsg{
		StateID:  pbft.NewStateID("c2"),
		SenderID: "s1",
	}

	// when
	err = c.SaveCommitMsg(cMsg)

	//then
	assert.Error(t, err)
	assert.Equal(t, 1, len(c.CommitMsgPool.Get()))
}
