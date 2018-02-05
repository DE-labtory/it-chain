package domain

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewConsensusState(t *testing.T) {
	viewID := "view"
	consensusID := "consensus"
	block := &Block{}

	consensusState := NewConsensusState(viewID,consensusID,block,PrePrepared)

	assert.Equal(t,consensusState.ViewID,viewID)
	assert.Equal(t,consensusState.ID,consensusID)
	assert.Equal(t,consensusState.CurrentStage,PrePrepared)
	assert.NotNil(t,consensusState.CommitMsgs)
	assert.NotNil(t,consensusState.PrepareMsgs)
}

func TestNewConsesnsusMessage(t *testing.T) {

	viewID := "view"
	block := &Block{}

	message:= NewConsesnsusMessage(viewID,1,block,"peer1",PrepareMsg)

	assert.Equal(t,message.ViewID,viewID)
	assert.Equal(t,message.SequenceID,int64(1))
	assert.Equal(t,message.MsgType,PrepareMsg)
	assert.Equal(t,message.Block,block)
}