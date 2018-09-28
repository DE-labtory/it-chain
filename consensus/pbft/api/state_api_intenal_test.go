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

func TestConsensusApi_StartConsensus_State(t *testing.T) {

	tests := map[string]struct {
		input struct {
			block                         pbft.ProposedBlock
			isNeedConsensus               bool
			peerNum                       int
			isPrevoteConditionSatisfied   bool
			isPreCommitConditionSatisfied bool
		}
		err   error
		stage pbft.Stage
	}{
		"Case : Consensus가 필요하고 Proposed된 Block이 정상인경우 (Normal Case)": {
			input: struct {
				block                         pbft.ProposedBlock
				isNeedConsensus               bool
				peerNum                       int
				isPrevoteConditionSatisfied   bool
				isPreCommitConditionSatisfied bool
			}{normalBlock, true, 5, false, false},
			err:   nil,
			stage: pbft.PROPOSE_STAGE,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, true, false, false)
		assert.EqualValues(t, test.err, cApi.StartConsensus(test.input.block))
		loadedState, _ := cApi.repo.Load()
		assert.Equal(t, string(test.stage), string(loadedState.CurrentStage))
	}
}

func TestConsensusApi_HandleProposeMsg_State(t *testing.T) {

	var validLeaderProposeMsg = pbft.ProposeMsg{
		StateID: pbft.StateID{
			ID: "state1",
		},
		SenderID:       "Leader",
		Representative: nil,
		ProposedBlock: pbft.ProposedBlock{
			Seal: make([]byte, 0),
			Body: make([]byte, 0),
		},
	}

	tests := map[string]struct {
		input struct {
			proposeMsg      pbft.ProposeMsg
			isNeedConsensus bool
			peerNum         int
			isRepoFull      bool
		}
		err   error
		stage pbft.Stage
	}{
		"Case PrePrepareMsg의 Sender id와 Request된 Leader id가 일치하며, repo가 full이 아닌경우 (Normal Case)": {
			input: struct {
				proposeMsg      pbft.ProposeMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{validLeaderProposeMsg, false, 5, false},
			err:   nil,
			stage: pbft.PREVOTE_STAGE,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, true, false, false)
		assert.EqualValues(t, test.err, cApi.HandleProposeMsg(test.input.proposeMsg))
		loadedState, _ := cApi.repo.Load()
		assert.Equal(t, string(test.stage), string(loadedState.CurrentStage))
	}
}

func TestConsensusApi_HandlePrevoteMsg_State(t *testing.T) {

	var validPrevoteMsg = pbft.PrevoteMsg{
		StateID:   pbft.StateID{"state"},
		SenderID:  "user1",
		BlockHash: []byte{1, 2, 3, 5},
	}

	tests := map[string]struct {
		input struct {
			prevoteMsg      pbft.PrevoteMsg
			isNeedConsensus bool
			peerNum         int
			isRepoFull      bool
		}
		err   error
		stage pbft.Stage
	}{
		"case 1 PrepareMsg의 Cid와 repo의 Cid가 같고, repo에 consensus가 저장된경우 (Normal Case)": {
			input: struct {
				prevoteMsg      pbft.PrevoteMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{validPrevoteMsg, false, 5, true},
			err:   nil,
			stage: pbft.PRECOMMIT_STAGE,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, true, true, false)
		assert.EqualValues(t, test.err, cApi.HandlePrevoteMsg(test.input.prevoteMsg))
		loadedState, _ := cApi.repo.Load()
		assert.Equal(t, string(test.stage), string(loadedState.CurrentStage))
	}

}
func TestConsensusApi_RepositoryClone(t *testing.T) {
	// stateApi1 에는 setUpApiCondition에 의해 repo가 set된 상황
	stateApi1 := setUpApiCondition(false, 5, true, false, false)
	// stateApi2 에는 stateApi1의 Repo가 주입된 상황
	stateApi2 := NewStateApi("publish2", pbft.PropagateService{}, nil, nil, stateApi1.repo)

	stateApi1.repo.Remove()
	_, err := stateApi2.repo.Load()

	assert.Equal(t, pbft.ErrEmptyRepo, err)

	newState := pbft.State{
		StateID:          pbft.StateID{"state1"},
		Representatives:  nil,
		Block:            normalBlock,
		CurrentStage:     pbft.PREVOTE_STAGE,
		PrevoteMsgPool:   pbft.PrevoteMsgPool{},
		PreCommitMsgPool: pbft.PreCommitMsgPool{},
	}
	stateApi1.repo.Save(newState)
	_, err2 := stateApi2.repo.Load()

	assert.Equal(t, nil, err2)

}

func setUpApiCondition(isNeedConsensus bool, peerNum int, isNormalBlock bool,
	isPrepareConditionSatisfied bool, isCommitConditionSatisfied bool) StateApiImpl {

	reps := make([]*pbft.Representative, 0)
	for i := 0; i < 6; i++ {
		reps = append(reps, &pbft.Representative{
			ID: "user",
		})
	}
	prevoteMsgPool := pbft.NewPrevoteMsgPool()
	for i := 0; i < 5; i++ {
		senderStr := "sender"
		senderStr += string(i)
		prevoteMsgPool.Save(&pbft.PrevoteMsg{
			StateID:   pbft.StateID{"state"},
			SenderID:  senderStr,
			BlockHash: []byte{1, 2, 3, 4},
		})
	}

	precommitMsgPool := pbft.NewPreCommitMsgPool()
	for i := 0; i < 5; i++ {
		senderStr := "sender"
		senderStr += string(i)
		precommitMsgPool.Save(&pbft.PreCommitMsg{
			StateID:  pbft.StateID{"state"},
			SenderID: senderStr,
		})
	}

	mockEventService := mock.EventService{}
	mockEventService.PublishFunc = func(topic string, event interface{}) error {
		return nil
	}

	propagateService := pbft.NewPropagateService(mockEventService)

	parliamentService := &mock.ParliamentService{}
	parliamentService.RequestPeerListFunc = func() ([]pbft.MemberID, error) {
		peerList := make([]pbft.MemberID, peerNum)
		for i := 0; i < peerNum; i++ {
			userStr := "user"
			userStr += string(peerNum)
			peerList = append(peerList, pbft.MemberID(userStr))
		}

		return peerList, nil
	}
	parliamentService.IsNeedConsensusFunc = func() bool {
		return isNeedConsensus
	}
	parliamentService.RequestLeaderFunc = func() (pbft.MemberID, error) {
		return "Leader", nil
	}

	eventService := mock.EventService{}
	eventService.PublishFunc = func(topic string, event interface{}) error {
		return nil
	}
	eventService.ConfirmBlockFunc = func(block pbft.ProposedBlock) error {
		return nil
	}

	repo := mem.NewStateRepository()
	if isPrepareConditionSatisfied && isNormalBlock {

		savedConsensus := pbft.State{
			StateID:          pbft.StateID{"state"},
			Representatives:  reps,
			Block:            normalBlock,
			CurrentStage:     pbft.IDLE_STAGE,
			PrevoteMsgPool:   prevoteMsgPool,
			PreCommitMsgPool: pbft.PreCommitMsgPool{},
		}
		repo.Save(savedConsensus)

	} else if isCommitConditionSatisfied && !isNormalBlock {
		savedConsensus := pbft.State{
			StateID:          pbft.StateID{"state"},
			Representatives:  reps,
			Block:            errorBlock,
			CurrentStage:     pbft.IDLE_STAGE,
			PrevoteMsgPool:   pbft.PrevoteMsgPool{},
			PreCommitMsgPool: precommitMsgPool,
		}
		repo.Save(savedConsensus)
	}
	cApi := NewStateApi("my", propagateService, eventService, parliamentService, &repo)

	return cApi
}
