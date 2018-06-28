package consensus

import (
	"errors"
	"fmt"

	"github.com/it-chain/midgard"
	"sync"
)

type ParliamentId string

func (pId ParliamentId) ToString() string {
	return string(pId)
}

type Parliament struct {
	ParliamentId ParliamentId
	Leader       *Leader
	Members      []*Member
}

func NewParliament() Parliament {
	return Parliament{
		ParliamentId: ParliamentId("0"),
		Members:      make([]*Member, 0),
		Leader:       nil,
	}
}

func (p *Parliament) findIndexOfMember(memberID string) int {
	for i, member := range p.Members {
		if member.MemberId.Id == memberID {
			return i
		}
	}

	return -1
}

func (p *Parliament) On(event midgard.Event) error {
	switch v := event.(type) {

	case *LeaderChangedEvent:
		p.Leader.LeaderId = LeaderId{v.GetID()}

	case *MemberJoinedEvent:
		p.Members = append(p.Members, &Member{MemberId: MemberId{v.GetID()}})

	case *MemberRemovedEvent:
		index := p.findIndexOfMember(v.GetID())

		if index != -1 {
			p.Members = append(p.Members[:index], p.Members[index+1:]...)
		}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

type ParliamentRepository interface {
	IsNeedConsensus() bool
	HasLeader() bool
	SetLeader(leader *Leader) midgard.Event
	GetLeader() *Leader
	AddMember(member *Member) (midgard.Event, error)
	RemoveMember(memberID MemberId) (midgard.Event, error)
	ValidateRepresentative(representatives []*Representative) bool
	FindByPeerID(memberID string) *Member
}

type ParliamentRepositoryImpl struct {
	lock       *sync.RWMutex
	parliament Parliament
}

func NewParliamentRepository() ParliamentRepository {
	return &ParliamentRepositoryImpl{
		lock:       &sync.RWMutex{},
		parliament: NewParliament(),
	}
}

func (pr *ParliamentRepositoryImpl) IsNeedConsensus() bool {
	numOfMember := 0

	if pr.HasLeader() {
		numOfMember = numOfMember + 1
	}

	numOfMember = numOfMember + len(pr.parliament.Members)

	if numOfMember >= 1 {
		return true
	}

	return false
}

func (pr *ParliamentRepositoryImpl) HasLeader() bool {
	if pr.parliament.Leader == nil {
		return false
	}

	return true
}

func (pr *ParliamentRepositoryImpl) SetLeader(leader *Leader) midgard.Event {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	pr.parliament.Leader = leader

	return LeaderChangedEvent{
		EventModel: midgard.EventModel{
			ID: leader.GetId(),
		},
	}
}

func (pr *ParliamentRepositoryImpl) GetLeader() *Leader {
	return pr.parliament.Leader
}

func (pr *ParliamentRepositoryImpl) AddMember(member *Member) (midgard.Event, error) {
	if member == nil {
		return nil, errors.New("Member is nil")
	}

	if member.GetId() == "" {
		return nil, errors.New(fmt.Sprintf("Need Valid PeerID [%s]", member.GetId()))
	}

	index := pr.parliament.findIndexOfMember(member.GetId())

	if index != -1 {
		return nil, errors.New(fmt.Sprintf("Already exist member [%s]", member.GetId()))
	}

	pr.parliament.Members = append(pr.parliament.Members, member)

	return MemberJoinedEvent{
		EventModel: midgard.EventModel{
			ID: member.GetId(),
		},
	}, nil
}

func (pr *ParliamentRepositoryImpl) RemoveMember(memberID MemberId) (midgard.Event, error) {

	index := pr.parliament.findIndexOfMember(memberID.ToString())

	if index == -1 {
		return nil, nil
	}

	pr.parliament.Members = append(pr.parliament.Members[:index], pr.parliament.Members[index+1:]...)

	return MemberRemovedEvent{
		EventModel: midgard.EventModel{
			ID: memberID.ToString(),
		},
	}, nil
}

func (pr *ParliamentRepositoryImpl) ValidateRepresentative(representatives []*Representative) bool {

	for _, representatives := range representatives {
		index := pr.parliament.findIndexOfMember(representatives.GetIdString())

		if index == -1 {
			return false
		}
	}

	return true
}

func (pr *ParliamentRepositoryImpl) FindByPeerID(memberID string) *Member {

	index := pr.parliament.findIndexOfMember(memberID)

	if index == -1 {
		return nil
	}

	return pr.parliament.Members[index]
}
