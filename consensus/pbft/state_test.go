/*
 * Copyright 2018 DE-labtory
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
	"fmt"
	"strconv"
	"testing"

	"github.com/DE-labtory/it-chain/consensus/pbft"
	"github.com/stretchr/testify/assert"
)

func TestPrevoteMsgPool_Save(t *testing.T) {
	// given
	pPool := pbft.NewPrevoteMsgPool()

	// case 1 : save
	pMsg := pbft.PrevoteMsg{
		MsgID:     "m1",
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 1, len(pPool.FindAll()))

	// case 2 : save
	pMsg = pbft.PrevoteMsg{
		MsgID:     "m2",
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s2",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.FindAll()))

	// case 3 : same sender
	pMsg = pbft.PrevoteMsg{
		MsgID:     "m3",
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s2",
		BlockHash: make([]byte, 0),
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.FindAll()))

	// case 4 : block hash is is nil
	pMsg = pbft.PrevoteMsg{
		MsgID:     "m4",
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s3",
		BlockHash: nil,
	}

	// when
	pPool.Save(&pMsg)

	// then
	assert.Equal(t, 2, len(pPool.FindAll()))
}

func TestPrevoteMsgPool_RemoveAllMsgs(t *testing.T) {
	// given
	pPool := pbft.NewPrevoteMsgPool()

	pMsg := pbft.PrevoteMsg{
		MsgID:     "m1",
		StateID:   pbft.StateID{"c1"},
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	pPool.Save(&pMsg)

	// when
	pPool.RemoveAllMsgs()

	// then
	assert.Equal(t, 0, len(pPool.FindAll()))
}

func TestPrevoteMsgPool_Remove(t *testing.T) {
	// given
	pPool := pbft.NewPrevoteMsgPool()

	for i := 1; i < 4; i++ {
		pMsg := pbft.PrevoteMsg{
			MsgID:     "m" + strconv.Itoa(i),
			StateID:   pbft.StateID{"state1"},
			SenderID:  "sender" + strconv.Itoa(i),
			BlockHash: []byte{1, 2, 3, 4},
		}
		pPool.Save(&pMsg)
	}

	assert.Equal(t, 3, len(pPool.FindAll()))

	// when : success
	pPool.Remove("m3")

	// then
	assert.Equal(t, 2, len(pPool.FindAll()))

	// when : fail
	pPool.Remove("m4")

	// then
	assert.Equal(t, 2, len(pPool.FindAll()))
}

func TestPrevoteMsgPool_FindAll(t *testing.T) {
	// given
	pPool := pbft.NewPrevoteMsgPool()

	// when
	numOfMsg := len(pPool.FindAll())

	// then
	assert.Equal(t, 0, numOfMsg)

	// given
	for i := 1; i < 4; i++ {
		pMsg := pbft.PrevoteMsg{
			MsgID:     "m" + strconv.Itoa(i),
			StateID:   pbft.StateID{"state1"},
			SenderID:  "sender" + strconv.Itoa(i),
			BlockHash: []byte{1, 2, 3, 4},
		}
		pPool.Save(&pMsg)
	}

	// when
	result := pPool.FindAll()
	numOfMsg = len(result)

	// then
	assert.Equal(t, 3, numOfMsg)
}

func TestPrevoteMsgPool_FindById(t *testing.T) {
	// given
	pPool := pbft.NewPrevoteMsgPool()

	for i := 1; i < 4; i++ {
		pMsg := pbft.PrevoteMsg{
			MsgID:     "m" + strconv.Itoa(i),
			StateID:   pbft.StateID{"state1"},
			SenderID:  "sender" + strconv.Itoa(i),
			BlockHash: []byte{1, 2, 3, 4},
		}
		pPool.Save(&pMsg)
	}

	// when : success
	result := pPool.FindById("m1")

	// then
	assert.Equal(t, "m1", result.MsgID)

	// when : fail
	result = pPool.FindById("m4")

	// then
	assert.Equal(t, "", result.MsgID)
}

func TestPreCommitMsgPool_Save(t *testing.T) {
	// given
	cPool := pbft.NewPreCommitMsgPool()

	// case 1 : save
	cMsg := pbft.PreCommitMsg{
		MsgID:    "m1",
		StateID:  pbft.StateID{"c1"},
		SenderID: "s1",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 1, len(cPool.FindAll()))

	// case 2 : save
	cMsg = pbft.PreCommitMsg{
		MsgID:    "m2",
		StateID:  pbft.StateID{"c1"},
		SenderID: "s2",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 2, len(cPool.FindAll()))

	// case 3 : same sender
	cMsg = pbft.PreCommitMsg{
		MsgID:    "m3",
		StateID:  pbft.StateID{"c1"},
		SenderID: "s2",
	}

	// when
	cPool.Save(&cMsg)

	// then
	assert.Equal(t, 2, len(cPool.FindAll()))
}

func TestPreCommitMsgPool_RemoveAllMsgs(t *testing.T) {
	// given
	cPool := pbft.NewPreCommitMsgPool()

	cMsg := pbft.PreCommitMsg{
		MsgID:    "m1",
		StateID:  pbft.StateID{"c1"},
		SenderID: "s1",
	}

	cPool.Save(&cMsg)

	// when
	cPool.RemoveAllMsgs()

	// then
	assert.Equal(t, 0, len(cPool.FindAll()))
}

func TestPreCommitMsgPool_Remove(t *testing.T) {
	// given
	pPool := pbft.NewPreCommitMsgPool()

	for i := 1; i < 4; i++ {
		pMsg := pbft.PreCommitMsg{
			MsgID:    "m" + strconv.Itoa(i),
			StateID:  pbft.StateID{"state1"},
			SenderID: "sender" + strconv.Itoa(i),
		}
		pPool.Save(&pMsg)
	}

	assert.Equal(t, 3, len(pPool.FindAll()))

	// when : success
	pPool.Remove("m3")

	// then
	assert.Equal(t, 2, len(pPool.FindAll()))

	// when : fail
	pPool.Remove("m4")

	// then
	assert.Equal(t, 2, len(pPool.FindAll()))
}

func TestPreCommiteMsgPool_FindAll(t *testing.T) {
	// given
	pPool := pbft.NewPreCommitMsgPool()

	// when
	numOfMsg := len(pPool.FindAll())

	// then
	assert.Equal(t, 0, numOfMsg)

	// given
	for i := 1; i < 4; i++ {
		pMsg := pbft.PreCommitMsg{
			MsgID:    "m" + strconv.Itoa(i),
			StateID:  pbft.StateID{"state1"},
			SenderID: "sender" + strconv.Itoa(i),
		}
		pPool.Save(&pMsg)
	}

	// when
	result := pPool.FindAll()
	numOfMsg = len(result)

	// then
	assert.Equal(t, 3, numOfMsg)
}

func TestPreCommitMsgPool_FindById(t *testing.T) {
	// given
	pPool := pbft.NewPreCommitMsgPool()

	for i := 1; i < 4; i++ {
		pMsg := pbft.PreCommitMsg{
			MsgID:    "m" + strconv.Itoa(i),
			StateID:  pbft.StateID{"state1"},
			SenderID: "sender" + strconv.Itoa(i),
		}
		pPool.Save(&pMsg)
	}

	// when : success
	result := pPool.FindById("m1")

	// then
	assert.Equal(t, "m1", result.MsgID)

	// when : fail
	result = pPool.FindById("m4")

	// then
	assert.Equal(t, "", result.MsgID)
}

func TestState_GetReceipients(t *testing.T) {
	// given
	initRepresentatives := make([]pbft.Representative, 0)

	for i := 0; i < 3; i++ {
		initRepresentatives = append(initRepresentatives, pbft.NewRepresentative(fmt.Sprint("s", i)))
	}

	c := pbft.State{
		StateID:         pbft.NewStateID("c1"),
		Representatives: initRepresentatives,
		Block: pbft.ProposedBlock{
			Seal: make([]byte, 0),
		},
		CurrentStage:     pbft.IDLE_STAGE,
		PrevoteMsgPool:   pbft.NewPrevoteMsgPool(),
		PreCommitMsgPool: pbft.NewPreCommitMsgPool(),
	}

	// case 1 : publisher is not in Representatives
	receipients := c.GetReceipients("s4")

	// then
	assert.Equal(t, c.Representatives, receipients)

	// case 2 : publisher is in Representatives
	receipients = c.GetReceipients("s2")

	correctReceipients := make([]pbft.Representative, 0)
	for i := 0; i < 2; i++ {
		correctReceipients = append(correctReceipients, pbft.NewRepresentative(fmt.Sprint("s", i)))
	}

	// then
	assert.Equal(t, correctReceipients, receipients)

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
		MsgID:     "m1",
		StateID:   pbft.NewStateID("c1"),
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	err := c.SavePrevoteMsg(pMsg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.PrevoteMsgPool.FindAll()))
	assert.Equal(t, *pMsg, c.PrevoteMsgPool.FindAll()[0])

	// case 2 : incorrect consensus ID
	pMsg = &pbft.PrevoteMsg{
		MsgID:     "m2",
		StateID:   pbft.NewStateID("c2"),
		SenderID:  "s1",
		BlockHash: make([]byte, 0),
	}

	// when
	err = c.SavePrevoteMsg(pMsg)

	//then
	assert.Error(t, err)
	assert.Equal(t, 1, len(c.PrevoteMsgPool.FindAll()))
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
		MsgID:    "m1",
		StateID:  pbft.NewStateID("c1"),
		SenderID: "s1",
	}

	// when
	err := c.SavePreCommitMsg(cMsg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.PreCommitMsgPool.FindAll()))
	assert.Equal(t, *cMsg, c.PreCommitMsgPool.FindAll()[0])

	// case 2 : incorrect consensus ID
	cMsg = &pbft.PreCommitMsg{
		MsgID:    "m2",
		StateID:  pbft.NewStateID("c2"),
		SenderID: "s1",
	}

	// when
	err = c.SavePreCommitMsg(cMsg)

	//then
	assert.Error(t, err)
	assert.Equal(t, 1, len(c.PreCommitMsgPool.FindAll()))
}
