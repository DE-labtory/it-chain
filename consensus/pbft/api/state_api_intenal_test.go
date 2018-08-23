package api

import (
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
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
			block                       pbft.ProposedBlock
			isNeedConsensus             bool
			peerNum                     int
			isPrepareConditionSatisfied bool
			isCommitConditionSatisfied  bool
		}
		err   error
		stage pbft.Stage
	}{
		"Case : Consensus가 필요하고 Proposed된 Block이 정상인경우 (Normal Case)": {
			input: struct {
				block                       pbft.ProposedBlock
				isNeedConsensus             bool
				peerNum                     int
				isPrepareConditionSatisfied bool
				isCommitConditionSatisfied  bool
			}{normalBlock, true, 5, false, false},
			err:   nil,
			stage: pbft.PREPREPARE_STAGE,
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

func TestConsensusApi_HandlePrePrepareMsg_State(t *testing.T) {

	var validLeaderPrePrepareMsg = pbft.PrePrepareMsg{
		StateID:        pbft.StateID{},
		SenderID:       "Leader",
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
		err   error
		stage pbft.Stage
	}{
		"Case PrePrepareMsg의 Sender id와 Request된 Leader id가 일치하며, repo가 full이 아닌경우 (Normal Case)": {
			input: struct {
				preprePareMsg   pbft.PrePrepareMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{validLeaderPrePrepareMsg, false, 5, false},
			err:   nil,
			stage: pbft.PREPARE_STAGE,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, true, false, false)
		assert.EqualValues(t, test.err, cApi.HandlePrePrepareMsg(test.input.preprePareMsg))
		loadedState, _ := cApi.repo.Load()
		assert.Equal(t, string(test.stage), string(loadedState.CurrentStage))
	}
}

func TestConsensusApi_HandlePrepareMsg_State(t *testing.T) {

	var validPrepareMsg = pbft.PrepareMsg{
		StateID:   pbft.StateID{"state"},
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
		err   error
		stage pbft.Stage
	}{
		"case 1 PrepareMsg의 Cid와 repo의 Cid가 같고, repo에 consensus가 저장된경우 (Normal Case)": {
			input: struct {
				prepareMsg      pbft.PrepareMsg
				isNeedConsensus bool
				peerNum         int
				isRepoFull      bool
			}{validPrepareMsg, false, 5, true},
			err:   nil,
			stage: pbft.COMMIT_STAGE,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s ", testName)
		cApi := setUpApiCondition(test.input.isNeedConsensus, test.input.peerNum, true, true, false)
		assert.EqualValues(t, test.err, cApi.HandlePrepareMsg(test.input.prepareMsg))
		loadedState, _ := cApi.repo.Load()
		assert.Equal(t, string(test.stage), string(loadedState.CurrentStage))
	}

}

func setUpApiCondition(isNeedConsensus bool, peerNum int, isNormalBlock bool,
	isPrepareConditionSatisfied bool, isCommitConditionSatisfied bool) StateApiImpl {

	reps := make([]*pbft.Representative, 0)
	for i := 0; i < 6; i++ {
		reps = append(reps, &pbft.Representative{
			ID: "user",
		})
	}
	prepareMsgPool := pbft.NewPrepareMsgPool()
	for i := 0; i < 5; i++ {
		senderStr := "sender"
		senderStr += string(i)
		prepareMsgPool.Save(&pbft.PrepareMsg{
			StateID:   pbft.StateID{"state"},
			SenderID:  senderStr,
			BlockHash: []byte{1, 2, 3, 4},
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

	eventService := mock.EventService{}
	eventService.PublishFunc = func(topic string, event interface{}) error {
		return nil
	}
	eventService.ConfirmBlockFunc = func(block pbft.ProposedBlock) error {
		return nil
	}

	repo := pbft.NewStateRepository()
	if isPrepareConditionSatisfied && isNormalBlock {

		savedConsensus := pbft.State{
			StateID:         pbft.StateID{"state"},
			Representatives: reps,
			Block:           normalBlock,
			CurrentStage:    pbft.IDLE_STAGE,
			PrepareMsgPool:  prepareMsgPool,
			CommitMsgPool:   pbft.CommitMsgPool{},
		}
		repo.Save(savedConsensus)

	} else if isCommitConditionSatisfied && !isNormalBlock {
		savedConsensus := pbft.State{
			StateID:         pbft.StateID{"state"},
			Representatives: reps,
			Block:           errorBlock,
			CurrentStage:    pbft.IDLE_STAGE,
			PrepareMsgPool:  pbft.PrepareMsgPool{},
			CommitMsgPool:   commitMsgPool,
		}
		repo.Save(savedConsensus)
	}
	cApi := NewStateApi("my", propagateService, eventService, parliamentService, repo)

	return cApi
}
