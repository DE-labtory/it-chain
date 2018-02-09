package domain

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

var endFlag = false

var mockHandler = func (consensusState ConsensusState){
	endFlag = true
}

func TestNewConsensusState(t *testing.T) {

	consensusID := "consensus"
	block := &Block{}
	view := View{}

	view.ID = "123"
	view.LeaderID = "1"
	view.PeerID = []string{"1","2","3"}

	consensusState := NewConsensusState(view,consensusID,block,PrePrepared,mockHandler,5)

	assert.Equal(t,consensusState.View,view)
	assert.Equal(t,consensusState.ID,consensusID)
	assert.Equal(t,consensusState.CurrentStage,PrePrepared)
	assert.NotNil(t,consensusState.CommitMsgs)
	assert.NotNil(t,consensusState.PrepareMsgs)
}

func TestNewConsesnsusMessage(t *testing.T) {

	consensusID := "1"
	viewID := "view"
	block := &Block{}

	message:= NewConsesnsusMessage(consensusID,viewID,1,block,"peer1",PrepareMsg)

	assert.Equal(t,message.ViewID,viewID)
	assert.Equal(t,message.SequenceID,int64(1))
	assert.Equal(t,message.MsgType,PrepareMsg)
	assert.Equal(t,message.Block,block)
}


func TestFromConsensusProtoMessage(t *testing.T) {

}

func TestConsensusState_start(t *testing.T) {
	consensusID := "consensus"
	block := &Block{}

	view := View{}
	view.ID = "123"
	view.LeaderID = "1"
	view.PeerID = []string{"1","2","3"}

	NewConsensusState(view,consensusID,block,PrePrepared,mockHandler,3)

	//var period float32 = 0.2
	time.Sleep(6*time.Second)

	assert.True(t,endFlag)

	endFlag = false
}

func TestConsensusState_start2(t *testing.T) {
	consensusID := "consensus"
	block := &Block{}

	view := View{}
	view.ID = "123"
	view.LeaderID = "1"
	view.PeerID = []string{"1","2","3"}

	NewConsensusState(view,consensusID,block,PrePrepared,mockHandler,6)

	//var period float32 = 0.2
	time.Sleep(3*time.Second)

	assert.False(t,endFlag)
	endFlag = false
}

func TestConsensusState_End(t *testing.T) {
	consensusID := "consensus"
	block := &Block{}
	view := View{}
	view.ID = "123"
	view.LeaderID = "1"
	view.PeerID = []string{"1","2","3"}

	cs := NewConsensusState(view,consensusID,block,PrePrepared,mockHandler,3)

	cs.End()

	time.Sleep(6*time.Second)
	assert.False(t,endFlag)
}