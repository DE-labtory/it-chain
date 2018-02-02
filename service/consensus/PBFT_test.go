package consensus

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/service/blockchain"
)

func TestNewConsensusState(t *testing.T) {
	viewID := "view"
	consensusID := "consensus"
	block := &blockchain.Block{}

	consensusState := NewConsensusState(viewID,consensusID,block,PrePrepared)

	assert.Equal(t,consensusState.ViewID,viewID)
	assert.Equal(t,consensusState.ID,consensusID)
	assert.Equal(t,consensusState.CurrentStage,PrePrepared)
	assert.NotNil(t,consensusState.CommitMsgs)
	assert.NotNil(t,consensusState.PrepareMsgs)
}

func TestNewConsesnsusMessage(t *testing.T) {

}