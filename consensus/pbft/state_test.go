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

func TestPrevoteMsgPool_Save(t *testing.T) {
	// given
	pPool := pbft.NewPrevoteMsgPool()

	// case 1 : save
	pMsg := pbft.PrevoteMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 1, len(pPool.Get()))

	// case 2 : save
	pMsg = pbft.PrevoteMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s2",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.Get()))

	// case 3 : same sender
	pMsg = pbft.PrevoteMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s2",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.Get()))

	// case 4 : block hash is is nil
	pMsg = pbft.PrevoteMsg{
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s3",
		BlockHash: nil,
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.Get()))
}

func TestPrevoteMsgPool_RemoveAllMsgs(t *testing.T) {
	// given
	pPool := pbft.NewPrevoteMsgPool()

	pMsg := pbft.PrevoteMsg{
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

func TestPreCommitMsgPool_Save(t *testing.T) {
	// given
	cPool := pbft.NewPreCommitMsgPool()

	// case 1 : save
	cMsg := pbft.PreCommitMsg{
		StateID:  pbft.StateID{"c1"},
		SenderID: "s1",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 1, len(cPool.Get()))

	// case 2 : save
	cMsg = pbft.PreCommitMsg{
		StateID:  pbft.StateID{"c1"},
		SenderID: "s2",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 2, len(cPool.Get()))

	// case 3 : same sender
	cMsg = pbft.PreCommitMsg{
		StateID:  pbft.StateID{"c1"},
		SenderID: "s2",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 2, len(cPool.Get()))
}

func TestPreCommitMsgPool_RemoveAllMsgs(t *testing.T) {
	// given
	cPool := pbft.NewPreCommitMsgPool()

	cMsg := pbft.PreCommitMsg{
		StateID:  pbft.StateID{"c1"},
		SenderID: "s1",
	}

	cPool.Save(&cMsg)

	// when
	cPool.RemoveAllMsgs()

	// then
	assert.Equal(t, 0, len(cPool.Get()))
}

func TestState_SavePrevoteMsg(t *testing.T) {
	// given
	c := pbft.State{
		StateID:         pbft.NewStateID("c1"),
		Representatives: nil,
		Block: pbft.ProposedBlock{
			Seal: make([]byte, 0),
		},
		CurrentStage:     pbft.IDLE_STAGE,
		PrevoteMsgPool:   pbft.NewPrevoteMsgPool(),
		PreCommitMsgPool: pbft.NewPreCommitMsgPool(),
	}

	// case 1 : save
	pMsg := &pbft.PrevoteMsg{
		StateID:   pbft.NewStateID("c1"),
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	err := c.SavePrevoteMsg(pMsg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.PrevoteMsgPool.Get()))
	assert.Equal(t, *pMsg, c.PrevoteMsgPool.Get()[0])

	// case 2 : incorrect consensus ID
	pMsg = &pbft.PrevoteMsg{
		StateID:   pbft.NewStateID("c2"),
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	err = c.SavePrevoteMsg(pMsg)

	//then
	assert.Error(t, err)
	assert.Equal(t, 1, len(c.PrevoteMsgPool.Get()))
}

func TestState_SavePreCommitMsg(t *testing.T) {
	// given
	c := pbft.State{
		StateID:         pbft.NewStateID("c1"),
		Representatives: nil,
		Block: pbft.ProposedBlock{
			Seal: make([]byte, 0),
		},
		CurrentStage:     pbft.IDLE_STAGE,
		PrevoteMsgPool:   pbft.NewPrevoteMsgPool(),
		PreCommitMsgPool: pbft.NewPreCommitMsgPool(),
	}

	// case 1 : save
	cMsg := &pbft.PreCommitMsg{
		StateID:  pbft.NewStateID("c1"),
		SenderID: "s1",
	}

	// when
	err := c.SavePreCommitMsg(cMsg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.PreCommitMsgPool.Get()))
	assert.Equal(t, *cMsg, c.PreCommitMsgPool.Get()[0])

	// case 2 : incorrect consensus ID
	cMsg = &pbft.PreCommitMsg{
		StateID:  pbft.NewStateID("c2"),
		SenderID: "s1",
	}

	// when
	err = c.SavePreCommitMsg(cMsg)

	//then
	assert.Error(t, err)
	assert.Equal(t, 1, len(c.PreCommitMsgPool.Get()))
}
