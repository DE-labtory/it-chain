package memory

import (
	"sync"

	"github.com/it-chain/it-chain-Engine/txpool"
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

func (lr *MemoryLeaderRepository) SetLeader(leader txpool.Leader) {
	lr.lock.Lock()
	defer lr.lock.Unlock()

	lr.currentLeader = leader
}
