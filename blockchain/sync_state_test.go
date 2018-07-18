package blockchain_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestBlockSyncState(t *testing.T) {
	// when
	syncState := blockchain.NewBlockSyncState()

	// then
	assert.Equal(t, blockchain.BC_SYNC_STATE_AID, syncState.GetID())

	// When
	event1 := &blockchain.SyncStartEvent{
		EventModel: midgard.EventModel{
			ID: blockchain.BC_SYNC_STATE_AID,
		},
	}
	syncState.On(event1)

	// Then
	assert.Equal(t, blockchain.PROGRESSING, syncState.IsProgressing())

	// When
	event2 := &blockchain.SyncDoneEvent{
		EventModel: midgard.EventModel{
			ID: blockchain.BC_SYNC_STATE_AID,
		},
	}
	syncState.On(event2)

	// Then
	assert.Equal(t, blockchain.DONE, syncState.IsProgressing())
}

func TestBlockSyncState_SetProgress(t *testing.T) {
	// when
	syncState := blockchain.NewBlockSyncState()

	// when
	syncState.SetProgress(blockchain.PROGRESSING)

	// then
	assert.Equal(t, blockchain.PROGRESSING, syncState.IsProgressing())

	// when
	syncState.SetProgress(blockchain.DONE)

	// then
	assert.Equal(t, blockchain.DONE, syncState.IsProgressing())
}
