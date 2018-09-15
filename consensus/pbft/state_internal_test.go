package pbft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// When Representative Number : 6, prevoteMsg Num : 4 -> then true
func TestState_CheckPrevoteCondition_Satisfy(t *testing.T) {
	// 6 rep
	satisfyPrevoteConditionState := setUpState()
	prevoteMsgs := make([]PrevoteMsg, 0)
	for i := 0; i < 4; i++ {
		prevoteMsgs = append(prevoteMsgs, PrevoteMsg{
			StateID:  StateID{"state1"},
			SenderID: "user1",
		})
	}
	satisfyPrevoteConditionState.PrevoteMsgPool = PrevoteMsgPool{messages: prevoteMsgs}

	assert.Equal(t, true, satisfyPrevoteConditionState.CheckPrevoteCondition())
}

// When Representative Number : 6, prevoteMsg Number : 3 -> then false
func TestState_CheckPrevoteCondition_UnSatisfy(t *testing.T) {
	unSatisfyPrevoteConditionState := setUpState()
	prevoteMsgs := make([]PrevoteMsg, 0)
	for i := 0; i < 3; i++ {
		prevoteMsgs = append(prevoteMsgs, PrevoteMsg{
			StateID:  StateID{"state1"},
			SenderID: "user1",
		})
	}
	unSatisfyPrevoteConditionState.PrevoteMsgPool = PrevoteMsgPool{messages: prevoteMsgs}

	//then false
	assert.Equal(t, false, unSatisfyPrevoteConditionState.CheckPrevoteCondition())
}

// When Representative Number : 6, prevoteCommitMsg Number : 4 -> then true
func TestState_CheckPreCommitCondition_Satisfy(t *testing.T) {
	satisfyPrecommitConditionState := setUpState()
	precommitMsgs := make([]PreCommitMsg, 0)
	for i := 0; i < 4; i++ {
		precommitMsgs = append(precommitMsgs, PreCommitMsg{
			StateID:  StateID{"state1"},
			SenderID: "user1",
		})
	}
	satisfyPrecommitConditionState.PreCommitMsgPool = PreCommitMsgPool{messages: precommitMsgs}

	assert.Equal(t, true, satisfyPrecommitConditionState.CheckPreCommitCondition())
}

// When Representative Number : 6, prevoteCommitMsg Number : 3 -> then false
func TestState_CheckPreCommitCondition_UnSatisfy(t *testing.T) {

	unSatisfyPrecommitConditionState := setUpState()
	precommitMsgs := make([]PreCommitMsg, 0)
	for i := 0; i < 3; i++ {
		precommitMsgs = append(precommitMsgs, PreCommitMsg{
			StateID:  StateID{"state1"},
			SenderID: "user1",
		})
	}
	unSatisfyPrecommitConditionState.PreCommitMsgPool = PreCommitMsgPool{messages: precommitMsgs}

	assert.Equal(t, false, unSatisfyPrecommitConditionState.CheckPreCommitCondition())
}

func setUpState() State {
	//when representatives consist 6 member
	reps := make([]*Representative, 0)
	for i := 0; i < 6; i++ {
		reps = append(reps, &Representative{
			ID: "user",
		})
	}

	var normalBlock = ProposedBlock{
		Seal: []byte{1, 2, 3, 4},
		Body: []byte{1, 2, 3, 5},
	}

	state1 := State{
		StateID:         StateID{"state1"},
		Representatives: reps,
		Block:           normalBlock,
		CurrentStage:    IDLE_STAGE,
	}
	return state1
}
