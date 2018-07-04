package blockchain_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/it-chain/midgard"
)

func TestBlockPoolModel(t *testing.T) {
	pool := *blockchain.NewBlockPool()

	block1 := &blockchain.DefaultBlock{
		Height: blockchain.BlockHeight(2),
	}

	// When
	pool.Add(block1)

	// Then
	assert.Equal(t, uint64(2), pool.Get(blockchain.BlockHeight(2)).GetHeight())


	// When

	pool.Delete(blockchain.BlockHeight(2))

	// Then
	assert.Equal(t, nil, pool.Get(blockchain.BlockHeight(2)))
}

func TestBlockSyncState(t *testing.T) {
	// when
	syncState := blockchain.NewBlockSyncState()

	// then
	assert.Equal(t, blockchain.BC_SYNC_STATE_AID, syncState.GetID())



	// When
	event1 := &blockchain.SyncStartEvent{
		EventModel: midgard.EventModel{
			ID: blockchain.BC_SYNC_STATE_EID,
		},
	}
	syncState.On(event1)

	// Then
	assert.Equal(t, blockchain.PROGRESSING, syncState.IsProgressing())


	// When
	event2 := &blockchain.SyncDoneEvent{
		EventModel: midgard.EventModel{
			ID: blockchain.BC_SYNC_STATE_EID,
		},
	}
	syncState.On(event2)

	// Then
	assert.Equal(t, blockchain.DONE, syncState.IsProgressing())
}
