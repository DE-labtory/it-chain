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

package api

import (
	"strconv"
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/infra/mem"
	"github.com/it-chain/engine/consensus/pbft/test/mock"
	"github.com/stretchr/testify/assert"
)

var normalBlock = pbft.ProposedBlock{
	Seal: []byte{1, 2, 3, 4},
	Body: []byte{1, 2, 3, 5},
}

var errorBlock = pbft.ProposedBlock{
	Seal: nil,
	Body: nil,
}

func TestStateApi_StartConsensus_State(t *testing.T) {

	tests := map[string]struct {
		input struct {
			block      pbft.ProposedBlock
			peerNum    int
			isRepoFull bool
		}
		err   error
		stage pbft.Stage
	}{
		"Case : Consensus가 필요하고 Proposed된 Block이 정상인경우 (Normal Case)": {
			input: struct {
				block      pbft.ProposedBlock
				peerNum    int
				isRepoFull bool
			}{normalBlock, 5, false},
			err:   nil,
			stage: pbft.PROPOSE_STAGE,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.peerNum, test.input.isRepoFull, true, pbft.IDLE_STAGE)
		cApi.StartConsensus(test.input.block)
		loadedState, _ := cApi.stateRepository.Load()
		assert.Equal(t, string(test.stage), string(loadedState.CurrentStage))
	}
}

func TestStateApi_HandleProposeMsg_CheckState(t *testing.T) {

	var validLeaderProposeMsg = pbft.ProposeMsg{
		StateID:        pbft.StateID{ID: "state"},
		SenderID:       "user0",
		Representative: nil,
		ProposedBlock: pbft.ProposedBlock{
			Seal: make([]byte, 0),
			Body: make([]byte, 0),
		},
	}

	tests := map[string]struct {
		input struct {
			proposeMsg pbft.ProposeMsg
			peerNum    int
			isRepoFull bool
		}
		err   error
		stage pbft.Stage
	}{
		"Case PrePrepareMsg의 Sender id와 Request된 Leader id가 일치하며, repo가 full이 아닌경우 (Normal Case)": {
			input: struct {
				proposeMsg pbft.ProposeMsg
				peerNum    int
				isRepoFull bool
			}{validLeaderProposeMsg, 5, false},
			err:   nil,
			stage: pbft.PREVOTE_STAGE,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.peerNum, false, false, pbft.IDLE_STAGE)
		assert.EqualValues(t, test.err, cApi.AcceptProposal(test.input.proposeMsg))
		loadedState, _ := cApi.stateRepository.Load()
		assert.Equal(t, string(test.stage), string(loadedState.CurrentStage))
	}
}

func TestStateApi_RepositoryClone(t *testing.T) {
	// stateApi1 에는 setUpApiCondition에 의해 repo가 set된 상황
	stateApi1 := setUpApiCondition(5, true, false, pbft.IDLE_STAGE)
	// stateApi2 에는 stateApi1의 Repo가 주입된 상황
	stateApi2 := NewStateApi("publish2", &pbft.PropagateService{}, nil, nil, stateApi1.stateRepository)

	stateApi1.stateRepository.Remove()
	_, err := stateApi2.stateRepository.Load()

	assert.Equal(t, pbft.ErrEmptyRepo, err)

	newState := pbft.State{
		StateID:          pbft.StateID{"state"},
		Representatives:  nil,
		Block:            normalBlock,
		CurrentStage:     pbft.PREVOTE_STAGE,
		PrevoteMsgPool:   pbft.PrevoteMsgPool{},
		PreCommitMsgPool: pbft.PreCommitMsgPool{},
	}
	stateApi1.stateRepository.Save(newState)
	_, err2 := stateApi2.stateRepository.Load()

	assert.Equal(t, nil, err2)

}

func TestStateApi_Reflect_TemporaryPrevoteMsgPool(t *testing.T) {

	reps := make([]pbft.Representative, 0)
	for i := 0; i < 5; i++ {
		reps = append(reps, pbft.Representative{
			ID: "user",
		})
	}
	var tempProposeMsg = pbft.ProposeMsg{
		StateID:        pbft.StateID{"state"},
		SenderID:       "user0",
		Representative: reps,
		ProposedBlock:  normalBlock,
	}

	var tempPrevoteMsg = pbft.PrevoteMsg{
		StateID:   pbft.StateID{"state"},
		SenderID:  "user1",
		BlockHash: []byte{1, 2, 3, 5},
	}

	var tempPrevoteMsg2 = pbft.PrevoteMsg{
		StateID:   pbft.StateID{"state"},
		SenderID:  "user2",
		BlockHash: []byte{1, 2, 3, 5},
	}

	//When Propose Msg를 받지못해 Repo에 State가 없음 then sApi의 tempPool이 저장 후 State가 생겼을 때 추가
	stateApi := setUpApiCondition(4, true, false, pbft.IDLE_STAGE)
	stateApi.stateRepository.Remove()

	stateApi.ReceivePrevote(tempPrevoteMsg)

	assert.Equal(t, 1, len(stateApi.tempPrevoteMsgPool.Get()))

	stateApi.AcceptProposal(tempProposeMsg)
	stateApi.ReceivePrevote(tempPrevoteMsg2)

	state, _ := stateApi.stateRepository.Load()
	assert.Equal(t, 2, len(state.PrevoteMsgPool.Get()))

}

func TestStateApi_Reflect_TemporaryPreCommitMsgPool(t *testing.T) {

	reps := make([]pbft.Representative, 0)
	for i := 0; i < 6; i++ {
		reps = append(reps, pbft.Representative{
			ID: "user",
		})
	}
	var tempProposeMsg = pbft.ProposeMsg{
		StateID:        pbft.StateID{"state"},
		SenderID:       "user0",
		Representative: reps,
		ProposedBlock:  normalBlock,
	}

	var tempPreCommitMsg = pbft.PreCommitMsg{
		StateID:  pbft.StateID{"state"},
		SenderID: "user1",
	}

	var tempPreCommitMsg2 = pbft.PreCommitMsg{
		StateID:  pbft.StateID{"state"},
		SenderID: "user2",
	}

	//When Propose Msg를 받지못해 Repo에 State가 없음 then sApi의 tempPool이 저장 후 State가 생겼을 때 추가
	stateApi := setUpApiCondition(6, true, false, pbft.IDLE_STAGE)
	stateApi.stateRepository.Remove()

	stateApi.ReceivePreCommit(tempPreCommitMsg)
	assert.Equal(t, 1, len(stateApi.tempPreCommitMsgPool.Get()))

	stateApi.AcceptProposal(tempProposeMsg)
	stateApi.ReceivePreCommit(tempPreCommitMsg2)
	state, _ := stateApi.stateRepository.Load()
	assert.Equal(t, 2, len(state.PreCommitMsgPool.Get()))

}

// todo
func TestStateApi_checkPrevote(t *testing.T) {

}

// todo
func TestStateApi_checkPreCommit(t *testing.T) {

}

func setUpApiCondition(peerNum int, isRepoFull, isNormalBlock bool, stage pbft.Stage) *StateApi {
	reps := make([]pbft.Representative, 0)
	for i := 0; i < peerNum; i++ {
		reps = append(reps, pbft.Representative{
			ID: "user",
		})
	}
	commitMsgPool := pbft.NewPreCommitMsgPool()

	for i := 0; i < peerNum; i++ {
		senderStr := "sender"
		senderStr += string(i)
		commitMsgPool.Save(&pbft.PreCommitMsg{
			StateID:  pbft.StateID{"state"},
			SenderID: senderStr,
		})
	}

	mockEventService := mock.EventService{}
	mockEventService.PublishFunc = func(topic string, event interface{}) error {
		return nil
	}

	propagateService := pbft.NewPropagateService(mockEventService)
	parliamentRepository := mem.NewParliamentRepository()
	parliament := pbft.NewParliament()
	for i := 0; i < peerNum; i++ {
		userStr := "user"
		userStr += strconv.Itoa(i)
		parliament.AddRepresentative(pbft.NewRepresentative(userStr))
	}

	parliament.SetLeader("user0")
	parliamentRepository.Save(parliament)

	eventService := mock.EventService{}
	eventService.PublishFunc = func(topic string, event interface{}) error {
		return nil
	}
	eventService.ConfirmBlockFunc = func(block pbft.ProposedBlock) error {
		return nil
	}

	repo := mem.NewStateRepository()
	if isNormalBlock {
		savedConsensus := pbft.State{
			StateID:          pbft.StateID{"state"},
			Representatives:  reps,
			Block:            normalBlock,
			CurrentStage:     stage,
			PrevoteMsgPool:   pbft.PrevoteMsgPool{},
			PreCommitMsgPool: pbft.PreCommitMsgPool{},
		}
		repo.Save(savedConsensus)
	} else {
		savedConsensus := pbft.State{
			StateID:          pbft.StateID{"state"},
			Representatives:  reps,
			Block:            errorBlock,
			CurrentStage:     stage,
			PrevoteMsgPool:   pbft.PrevoteMsgPool{},
			PreCommitMsgPool: pbft.PreCommitMsgPool{},
		}
		repo.Save(savedConsensus)
	}
	if !isRepoFull {
		repo.Remove()
	}

	cApi := NewStateApi("my", propagateService, eventService, parliamentRepository, repo)
	return cApi
}
