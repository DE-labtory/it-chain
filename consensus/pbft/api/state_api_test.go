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

	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/engine/consensus/api"
	"github.com/it-chain/engine/consensus/infra/mem"
	"github.com/it-chain/engine/consensus/test/mock"
	"github.com/stretchr/testify/assert"
)

var normalBlock = consensus.ProposedBlock{
	Seal: []byte{1, 2, 3, 4},
	Body: []byte{1, 2, 3, 5},
}

var errorBlock = consensus.ProposedBlock{
	Seal: nil,
	Body: nil,
}

func TestConsensusApi_StartConsensus(t *testing.T) {

	//case 1 Consensus가 필요없고, Proposed된 Block이 정상인 경우
	cApi1 := setUpConsensusCondition(false, 1, false)
	err1 := cApi1.StartConsensus("my", normalBlock)
	assert.Nil(t, err1)

	//case 1-1 Consensus가 필요없고, Proposed된 Block의 hash값이 nil인 경우
	err2 := cApi1.StartConsensus("user1", errorBlock)
	assert.Equal(t, errors.New("Block hash is nil"), err2)

	//case 2 Consensus가 필요하고, Repo가 차있지 않는 상태 (Normal Case)
	cApi2 := setUpConsensusCondition(true, 5, false)
	err3 := cApi2.StartConsensus("my", normalBlock)
	assert.Equal(t, nil, err3)

	//case 2-1 Consensus가 필요하고, Repo가 full인 우
	cApi3 := setUpConsensusCondition(true, 5, true)
	err4 := cApi3.StartConsensus("my", normalBlock)
	assert.Equal(t, mem.ConsensusAlreadyExistError, err4)

}

func TestConsensusApi_ReceivePrePrepareMsg(t *testing.T) {

	var validLeaderPrePrepareMsg = consensus.PrePrepareMsg{
		ConsensusId:    consensus.ConsensusId{},
		SenderId:       "Leader",
		Representative: nil,
		ProposedBlock:  consensus.ProposedBlock{},
	}
	var invalidLeaderPrePrepareMsg = consensus.PrePrepareMsg{
		ConsensusId:    consensus.ConsensusId{},
		SenderId:       "NoLeader",
		Representative: nil,
		ProposedBlock:  consensus.ProposedBlock{},
	}
	//case 1 PrePrepareMsg의 Sender id와 Request된 Leader id가 일치하며, repo가 full이 아닌경우 (Normal Case)
	cApi1 := setUpConsensusCondition(false, 5, false)
	err := cApi1.ReceivePrePrepareMsg(validLeaderPrePrepareMsg)
	assert.Equal(t, nil, err)

	//case 1-1 PrePrepareMsg의 Sender id와 Request된 Leader id가 일치하며, repo가 full인 경우
	cApi2 := setUpConsensusCondition(false, 5, true)
	err2 := cApi2.ReceivePrePrepareMsg(validLeaderPrePrepareMsg)
	assert.Equal(t, mem.ConsensusAlreadyExistError, err2)

	//case 2 PrePrepareMsg의 Sender id와 Request된 Leader Id가 일치하지 않을 경우
	cApi3 := setUpConsensusCondition(false, 5, false)
	err3 := cApi3.ReceivePrePrepareMsg(invalidLeaderPrePrepareMsg)
	assert.Equal(t, err3, consensus.InvalidLeaderIdError)

}

func TestConsensusApi_ReceivePrepareMsg(t *testing.T) {

	var validPrepareMsg = consensus.PrepareMsg{
		ConsensusId: consensus.ConsensusId{"consensus"},
		SenderId:    "user1",
		BlockHash:   []byte{1, 2, 3, 5},
	}
	var invalidPrepareMsg = consensus.PrepareMsg{
		ConsensusId: consensus.ConsensusId{"invalidConsensus"},
		SenderId:    "user1",
		BlockHash:   []byte{1, 2, 3, 5},
	}

	//case 1 PrepareMsg의 Cid와 repo의 Cid가 같고, repo에 consensus가 저장된경우 (Normal Case)
	cApi1 := setUpConsensusCondition(false, 5, true)
	err1 := cApi1.ReceivePrepareMsg(validPrepareMsg)
	assert.Equal(t, nil, err1)
	//case 2 PrepareMsg의 Cid와 repo의 Cid가 다를 경우
	cApi2 := setUpConsensusCondition(false, 5, true)
	err2 := cApi2.ReceivePrepareMsg(invalidPrepareMsg)
	assert.Equal(t, errors.New("Consensus ID is not same"), err2)
	//case 3 Repo에 Consensus가 저장되어있지 않은 경우
	cApi3 := setUpConsensusCondition(false, 5, false)
	err3 := cApi3.ReceivePrepareMsg(validPrepareMsg)
	assert.Equal(t, mem.LoadConsensusError, err3)

}

func TestConsensusApi_ReceiveCommitMsg(t *testing.T) {

	var validCommitMsg = consensus.CommitMsg{
		ConsensusId: consensus.ConsensusId{"consensus"},
		SenderId:    "user1",
	}
	var invalidCommitMsg = consensus.CommitMsg{
		ConsensusId: consensus.ConsensusId{"invalidConsensus"},
		SenderId:    "user2",
	}
	//case 1 repo에 consensus가 있고 repo의 cid와 commitMsg의 cid가 일치하는 경우(Normal Case)
	cApi1 := setUpConsensusCondition(false, 5, true)
	err1 := cApi1.ReceiveCommitMsg(validCommitMsg)
	assert.Equal(t, nil, err1)

	//case 2 repo에 consensus가 저장되어있지 않은 경우
	cApi2 := setUpConsensusCondition(false, 5, false)
	err2 := cApi2.ReceiveCommitMsg(validCommitMsg)
	assert.Equal(t, mem.LoadConsensusError, err2)

	//case 3 repo에 저장된 consensus의 cid와 commitMsg의 cid가 일치하지 않은 경우
	cApi3 := setUpConsensusCondition(false, 5, true)
	err3 := cApi3.ReceiveCommitMsg(invalidCommitMsg)
	assert.Equal(t, errors.New("Consensus ID is not same"), err3)
}

func setUpConsensusCondition(isNeedConsensus bool, peerNum int, isRepoFull bool) api.ConsensusApi {
	propagateService := &mock.MockPropagateService{}
	propagateService.BroadcastPrePrepareMsgFunc = func(msg consensus.PrePrepareMsg) error {
		return nil
	}
	propagateService.BroadcastPrepareMsgFunc = func(msg consensus.PrepareMsg) error {
		return nil
	}
	propagateService.BroadcastCommitMsgFunc = func(msg consensus.CommitMsg) error {
		return nil
	}

	parliamentService := &mock.MockParliamentService{}
	parliamentService.RequestPeerListFunc = func() ([]consensus.MemberId, error) {
		peerList := make([]consensus.MemberId, peerNum)
		for i := 0; i < peerNum; i++ {
			userStr := "user"
			userStr += string(peerNum)
			peerList = append(peerList, consensus.MemberId(userStr))
		}

		return peerList, nil
	}
	parliamentService.IsNeedConsensusFunc = func() bool {
		return isNeedConsensus
	}
	parliamentService.RequestLeaderFunc = func() (consensus.MemberId, error) {
		return "Leader", nil
	}

	confirmService := &mock.MockConfirmService{}
	confirmService.ConfirmBlockFunc = func(block consensus.ProposedBlock) error {
		if block.Seal == nil {
			return errors.New("Block hash is nil")
		}

		if block.Body == nil {
			return errors.New("There is no block")
		}
		return nil
	}
	repo := mem.NewConsensusRepository()
	if isRepoFull {
		savedConsensus := consensus.Consensus{
			ConsensusID:     consensus.ConsensusId{"consensus"},
			Representatives: nil,
			Block:           normalBlock,
			CurrentState:    "",
			PrepareMsgPool:  consensus.PrepareMsgPool{},
			CommitMsgPool:   consensus.CommitMsgPool{},
		}
		repo.Save(savedConsensus)
	}
	cApi := api.NewConsensusApi(propagateService, confirmService, parliamentService, repo)
	return cApi
}
