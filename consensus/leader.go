package consensus

import (
	"github.com/it-chain/midgard"
	"errors"
	"fmt"
	"sync"
)

type LeaderId struct {
	Id string
}

type Leader struct {
	LeaderId LeaderId
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

func (l *Leader) StringLeaderId() string {
	return l.LeaderId.ToString()
}

func (l Leader) GetID() string {
	return l.StringLeaderId()
}

func (l *Leader) On(event midgard.Event) error {
	switch v := event.(type) {

	case *LeaderChangedEvent:
		l.LeaderId = LeaderId{v.GetID()}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
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
