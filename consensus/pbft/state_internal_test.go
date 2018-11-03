package pbft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// When Representative Number : 6, prevoteMsg Num : 4 -> then true
func TestState_CheckPrevoteCondition_Satisfy(t *testing.T) {
	// 6 rep
	satisfyPrevoteConditionState := setUpState()
	satisfyPrevoteConditionState.PrevoteMsgPool = NewPrevoteMsgPool()
	for i := 0; i < 4; i++ {
		msg := PrevoteMsg{
			MsgID:     "msg" + string(i),
			StateID:   StateID{"state1"},
			SenderID:  "user" + string(i),
			BlockHash: []byte{1, 2, 3, 4},
		}
		satisfyPrevoteConditionState.PrevoteMsgPool.Save(&msg)
	}

	assert.Equal(t, true, satisfyPrevoteConditionState.CheckPrevoteCondition())
}

// When Representative Number : 6, prevoteMsg Number : 2 -> then false
func TestState_CheckPrevoteCondition_UnSatisfy(t *testing.T) {
	unSatisfyPrevoteConditionState := setUpState()
	unSatisfyPrevoteConditionState.PrevoteMsgPool = NewPrevoteMsgPool()
	for i := 0; i < 2; i++ {
		msg := PrevoteMsg{
			MsgID:     "msg" + string(i),
			StateID:   StateID{"state1"},
			SenderID:  "user" + string(i),
			BlockHash: []byte{1, 2, 3, 4},
		}
		unSatisfyPrevoteConditionState.PrevoteMsgPool.Save(&msg)

	}

	//then false
	assert.Equal(t, false, unSatisfyPrevoteConditionState.CheckPrevoteCondition())
}

// When Representative Number : 6, prevoteCommitMsg Number : 4 -> then true
func TestState_CheckPreCommitCondition_Satisfy(t *testing.T) {
	satisfyPrecommitConditionState := setUpState()
	satisfyPrecommitConditionState.PreCommitMsgPool = NewPreCommitMsgPool()
	for i := 0; i < 4; i++ {
		msg := PreCommitMsg{
			MsgID:    "msg" + string(i),
			StateID:  StateID{"state1"},
			SenderID: "user" + string(i),
		}
		satisfyPrecommitConditionState.PreCommitMsgPool.Save(&msg)
	}

	assert.Equal(t, true, satisfyPrecommitConditionState.CheckPreCommitCondition())
}

// When Representative Number : 6, prevoteCommitMsg Number : 2 -> then false
func TestState_CheckPreCommitCondition_UnSatisfy(t *testing.T) {
	unSatisfyPrecommitConditionState := setUpState()
	unSatisfyPrecommitConditionState.PreCommitMsgPool = NewPreCommitMsgPool()
	for i := 0; i < 2; i++ {
		msg := PreCommitMsg{
			MsgID:    "msg" + string(i),
			StateID:  StateID{"state1"},
			SenderID: "user" + string(i),
		}
		unSatisfyPrecommitConditionState.PreCommitMsgPool.Save(&msg)
	}

	assert.Equal(t, false, unSatisfyPrecommitConditionState.CheckPreCommitCondition())
}

func setUpState() State {
	//when representatives consist 6 member
	reps := make([]Representative, 0)
	for i := 0; i < 6; i++ {
		reps = append(reps, Representative{
			ID: "user" + string(i),
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
