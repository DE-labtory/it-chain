package txpool

import "sync"

type Leader struct {
	leaderId LeaderId
}

func (l *Leader) StringLeaderId() string {
	return l.leaderId.ToString()
}

type LeaderId struct {
	id string
}

func (lid LeaderId) ToString() string {
	return string(lid.id)
}

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
