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

package api_test

import (
	"testing"

	"errors"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
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

func TestConsensusApi_StartConsensus(t *testing.T) {

	tests := map[string]struct {
		input struct {
			block           pbft.ProposedBlock
			isNeedConsensus bool
			peerNum         int
			isRepoFull      bool
		}
		err error
	}{
		"Case1 : Consensus가 필요없는 상황": {
			input: struct {
				block           pbft.ProposedBlock
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{normalBlock, false, 1, false},
			err: api.ConsensusCreateError,
		},
		"Case2 : Consensus가 필요하고 Proposed된 Block이 정상이며, Repo가 차있지 않는 경우": {
			input: struct {
				block           pbft.ProposedBlock
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{normalBlock, true, 5, false},
			err: nil,
		},
		"Case3 : Consensus가 필요하고 repo가 Full인 경우": {
			input: struct {
				block           pbft.ProposedBlock
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{normalBlock, true, 5, true},
			err: mem.ErrConsensusAlreadyExist,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, test.input.isRepoFull, true)
		assert.EqualValues(t, test.err, cApi.StartConsensus(test.input.block))

	}

}
func TestConsensusApi_HandlePrePrepareMsg(t *testing.T) {

	var validLeaderPrePrepareMsg = pbft.PrePrepareMsg{
		StateID:        pbft.StateID{},
		SenderID:       "Leader",
		Representative: nil,
		ProposedBlock:  pbft.ProposedBlock{},
	}
	var invalidLeaderPrePrepareMsg = pbft.PrePrepareMsg{
		StateID:        pbft.StateID{},
		SenderID:       "NoLeader",
		Representative: nil,
		ProposedBlock:  pbft.ProposedBlock{},
	}
	tests := map[string]struct {
		input struct {
			preprePareMsg   pbft.PrePrepareMsg
			isNeedConsensus bool
			peerNum         int
			isRepoFull      bool
		}
		err error
	}{
		"Case 1 PrePrepareMsg의 Sender id와 Request된 Leader id가 일치하며, repo가 full이 아닌경우 (Normal Case)": {
			input: struct {
				preprePareMsg   pbft.PrePrepareMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{validLeaderPrePrepareMsg, false, 5, false},
			err: nil,
		},
		"Case 2 PrePrepareMsg의 Sender id와 Request된 Leader id가 일치하며, repo가 full인 경우": {
			input: struct {
				preprePareMsg   pbft.PrePrepareMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{validLeaderPrePrepareMsg, false, 5, true},
			err: mem.ErrConsensusAlreadyExist,
		},
		"Case 3 PrePrepareMsg의 Sender id와 Request된 Leader id가 일치하지 않을 경우": {
			input: struct {
				preprePareMsg   pbft.PrePrepareMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{invalidLeaderPrePrepareMsg, false, 5, false},
			err: pbft.InvalidLeaderIdError,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, test.input.isRepoFull, true)
		assert.EqualValues(t, test.err, cApi.HandlePrePrepareMsg(test.input.preprePareMsg))

	}
}

func TestConsensusApi_HandlePrepareMsg(t *testing.T) {

	var validPrepareMsg = pbft.PrepareMsg{
		StateID:   pbft.StateID{"state"},
		SenderID:  "user1",
		BlockHash: []byte{1, 2, 3, 5},
	}
	var invalidPrepareMsg = pbft.PrepareMsg{
		StateID:   pbft.StateID{"invalidState"},
		SenderID:  "user1",
		BlockHash: []byte{1, 2, 3, 5},
	}

	tests := map[string]struct {
		input struct {
			prepareMsg      pbft.PrepareMsg
			isNeedConsensus bool
			peerNum         int
			isRepoFull      bool
		}
		err error
	}{
		"case 1 PrepareMsg의 Cid와 repo의 Cid가 같고, repo에 consensus가 저장된경우 (Normal Case)": {
			input: struct {
				prepareMsg      pbft.PrepareMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{validPrepareMsg, false, 5, true},
			err: nil,
		},
		"Case 2 PrepareMsg의 Cid와 repo의 Cid가 다를 경우": {
			input: struct {
				prepareMsg      pbft.PrepareMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{invalidPrepareMsg, false, 5, true},
			err: errors.New("State ID is not same"),
		},
		"Case 3 Repo에 Consensus가 저장되어있지 않은 경우": {
			input: struct {
				prepareMsg      pbft.PrepareMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{validPrepareMsg, false, 5, false},
			err: mem.ErrLoadConsensus,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, test.input.isRepoFull, true)
		assert.EqualValues(t, test.err, cApi.HandlePrepareMsg(test.input.prepareMsg))
	}

}

func TestConsensusApi_HandleCommitMsg(t *testing.T) {

	var validCommitMsg = pbft.CommitMsg{
		StateID:  pbft.StateID{"state"},
		SenderID: "user1",
	}
	var invalidCommitMsg = pbft.CommitMsg{
		StateID:  pbft.StateID{"invalidState"},
		SenderID: "user2",
	}

	tests := map[string]struct {
		input struct {
			commitMsg       pbft.CommitMsg
			isNeedConsensus bool
			peerNum         int
			isRepoFull      bool
			isNormalBlock   bool
		}
		err error
	}{
		"Case 1 repo에 consensus가 있고 repo의 cid와 commitMsg의 cid가 일치하는 경우(Normal Case)": {
			input: struct {
				commitMsg       pbft.CommitMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
				isNormalBlock   bool
			}{validCommitMsg, false, 5, true, true},
			err: nil,
		},
		"Case 2 repo에 consensus가 저장되어있지 않은 경우": {
			input: struct {
				commitMsg       pbft.CommitMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
				isNormalBlock   bool
			}{validCommitMsg, false, 5, false, true},
			err: mem.ErrLoadConsensus,
		},
		"Case 3 repo에 저장된 pbft의 cid와 commitMsg의 cid가 일치하지 않은 경우": {
			input: struct {
				commitMsg       pbft.CommitMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
				isNormalBlock   bool
			}{invalidCommitMsg, false, 5, true, true},
			err: errors.New("State ID is not same"),
		},
		"Case 4 repo에 저장된 pbft cid와 commitMsg의 cid가 일치하고, Commit조건을 만족할 경우": {
			input: struct {
				commitMsg       pbft.CommitMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
				isNormalBlock   bool
			}{validCommitMsg, false, 5, true, false},
			err: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, test.input.isRepoFull, test.input.isNormalBlock)
		assert.EqualValues(t, test.err, cApi.HandleCommitMsg(test.input.commitMsg))
	}
}

func setUpApiCondition(isNeedConsensus bool, peerNum int, isRepoFull bool, isNormalBlock bool) api.StateApi {

	reps := make([]*pbft.Representative, 0)
	for i := 0; i < 6; i++ {
		reps = append(reps, &pbft.Representative{
			ID: "user",
		})
	}

	commitMsgPool := pbft.NewCommitMsgPool()
	for i := 0; i < 5; i++ {
		senderStr := "sender"
		senderStr += string(i)
		commitMsgPool.Save(&pbft.CommitMsg{
			StateID:  pbft.StateID{"state"},
			SenderID: senderStr,
		})
	}

	propagateService := &mock.MockPropagateService{}
	propagateService.BroadcastPrePrepareMsgFunc = func(msg pbft.PrePrepareMsg) error {
		return nil
	}
	propagateService.BroadcastPrepareMsgFunc = func(msg pbft.PrepareMsg) error {
		return nil
	}
	propagateService.BroadcastCommitMsgFunc = func(msg pbft.CommitMsg) error {
		return nil
	}

	parliamentService := &mock.MockParliamentService{}
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

	eventService := adapter.NewEventService(func(topic string, data interface{}) (err error) {
		return nil
	})

	repo := mem.NewStateRepository()
	if isRepoFull && isNormalBlock {

		savedConsensus := pbft.State{
			StateID:         pbft.StateID{"state"},
			Representatives: reps,
			Block:           normalBlock,
			CurrentStage:    "",
			PrepareMsgPool:  pbft.PrepareMsgPool{},
			CommitMsgPool:   pbft.CommitMsgPool{},
		}
		repo.Save(savedConsensus)

	} else if isRepoFull && !isNormalBlock {
		savedConsensus := pbft.State{
			StateID:         pbft.StateID{"state"},
			Representatives: reps,
			Block:           errorBlock,
			CurrentStage:    "",
			PrepareMsgPool:  pbft.PrepareMsgPool{},
			CommitMsgPool:   commitMsgPool,
		}
		repo.Save(savedConsensus)
	}
	cApi := api.NewStateApi("my", propagateService, eventService, parliamentService, repo)
	return cApi
}
