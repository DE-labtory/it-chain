package txpool

import (
	"sync"
)

type LeaderRepository interface {
	GetLeader() Leader
	SetLeader(leader Leader)
}

type LeaderRepositoryImpl struct {
	lock          *sync.RWMutex
	currentLeader Leader
}

func NewLeaderRepository() LeaderRepository {
	return &LeaderRepositoryImpl{
		lock:          &sync.RWMutex{},
		currentLeader: Leader{},
	}
}

func (lr *LeaderRepositoryImpl) GetLeader() Leader {
	return lr.currentLeader
}

func (lr *LeaderRepositoryImpl) SetLeader(leader Leader) {
	lr.lock.Lock()
	defer lr.lock.Unlock()

	lr.currentLeader = leader
}
