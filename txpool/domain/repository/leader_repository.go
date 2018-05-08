package repository

import (
	"sync"

	"github.com/it-chain/it-chain-Engine/txpool/domain/model"
)

type LeaderRepository interface {
	GetLeader() model.Leader
	SetLeader(leader model.Leader)
}

type LeaderRepositoryImpl struct {
	lock          *sync.RWMutex
	currentLeader model.Leader
}

func NewLeaderRepository() LeaderRepository {
	return &LeaderRepositoryImpl{
		lock:          &sync.RWMutex{},
		currentLeader: model.Leader{},
	}
}

func (lr *LeaderRepositoryImpl) GetLeader() model.Leader {
	return lr.currentLeader
}

func (lr *LeaderRepositoryImpl) SetLeader(leader model.Leader) {
	lr.lock.Lock()
	defer lr.lock.Unlock()

	lr.currentLeader = leader
}
