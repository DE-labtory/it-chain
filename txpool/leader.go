package txpool

import (
	"errors"
	"fmt"
	"sync"

	"github.com/it-chain/midgard"
)

type LeaderId struct {
	Id string
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

//Aggregate root must implement aggregate interface
type Leader struct {
	LeaderId LeaderId
}

// must implement id method
func (l Leader) GetID() string {
	return l.StringLeaderId()
}

// must implement on method
func (l *Leader) On(event midgard.Event) error {

	switch v := event.(type) {

	case *LeaderChangedEvent:
		l.LeaderId = LeaderId{v.GetID()}

	default:
		return errors.New(fmt.Sprintf("unhandled event_store [%s]", v))
	}

	return nil
}

func (l *Leader) StringLeaderId() string {
	return l.LeaderId.ToString()
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
