package txpool

import "sync"

type LeaderId struct {
	id string
}

func (lid LeaderId) ToString() string {
	return string(lid.id)
}

//Aggregate root must implement aggregate interface
type Leader struct {
	leaderId LeaderId
}

func (l Leader) GetID() string {
	return l.StringLeaderId()
}

func (l *Leader) StringLeaderId() string {
	return l.leaderId.ToString()
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
