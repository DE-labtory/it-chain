package memory

import (
	"sync"

	"github.com/it-chain/engine/txpool"
)

type MemoryLeaderRepository struct {
	lock          *sync.RWMutex
	currentLeader txpool.Leader
}

func NewLeaderRepository() MemoryLeaderRepository {
	return MemoryLeaderRepository{
		lock:          &sync.RWMutex{},
		currentLeader: txpool.Leader{},
	}
}

func (lr *MemoryLeaderRepository) GetLeader() txpool.Leader {
	lr.lock.Lock()
	defer lr.lock.Unlock()

	return lr.currentLeader
}


