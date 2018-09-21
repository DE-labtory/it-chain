package mem

import (
	"sync"

	"github.com/it-chain/engine/blockchain"
)

type SyncStateRepository struct {
	m     *sync.RWMutex
	State blockchain.SyncState
}

func NewSyncStateRepository() *SyncStateRepository {
	return &SyncStateRepository{
		m:     &sync.RWMutex{},
		State: blockchain.NewSyncState(false, false),
	}
}

func (r *SyncStateRepository) Get() blockchain.SyncState {
	r.m.Lock()
	defer r.m.Unlock()

	return r.State
}

func (r *SyncStateRepository) Set(state blockchain.SyncState) {
	r.m.Lock()
	defer r.m.Unlock()

	r.State = state
}
