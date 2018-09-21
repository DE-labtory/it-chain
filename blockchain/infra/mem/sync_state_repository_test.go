package mem_test

import (
	"testing"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/magiconair/properties/assert"
)

func TestNewSyncStateRepository(t *testing.T) {
	repo := mem.NewSyncStateRepository()

	assert.Equal(t, repo.State.SyncProgress(), false)
	assert.Equal(t, repo.State.Sync(), false)

	state := blockchain.NewSyncState(true, false)
	repo.Set(state)

	assert.Equal(t, repo.State.SyncProgress(), true)
	assert.Equal(t, repo.State.Sync(), false)

	state2 := repo.Get()

	assert.Equal(t, state2.SyncProgress(), true)
	assert.Equal(t, state2.Sync(), false)
}
