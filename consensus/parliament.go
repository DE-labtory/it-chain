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

func (p *Parliament) GetID() string {
	return p.ParliamentId.ToString()
}

func (p *Parliament) IsNeedConsensus() bool {
	numOfMember := 0

	if p.HasLeader() {
		numOfMember = numOfMember + 1
	}

	numOfMember = numOfMember + len(p.Members)

	if numOfMember >= 1 {
		return true
	}

	return false
}

func (p *Parliament) HasLeader() bool {
	if p.Leader == nil {
		return false
	}

	return true
}

func (p *Parliament) ChangeLeader(leader *Leader) (*LeaderChangedEvent, error) {
	if leader == nil {
		return nil, errors.New("Leader is nil")
	}

	leaderChangedEvent := LeaderChangedEvent{
		EventModel: midgard.EventModel{
			ID: p.GetID(),
		},
		LeaderId: leader.LeaderId,
	}

	p.On(&leaderChangedEvent)

	return &leaderChangedEvent, nil
}

func (p *Parliament) AddMember(member *Member) (*MemberJoinedEvent, error) {
	if member == nil {
		return nil, errors.New("Member is nil")
	}

	if member.GetId() == "" {
		return nil, errors.New(fmt.Sprintf("Need Valid PeerID [%s]", member.GetId()))
	}

	index := p.findIndexOfMember(member.GetId())

	if index != -1 {
		return nil, errors.New(fmt.Sprintf("Already exist member [%s]", member.GetId()))
	}

	memberJoinedEvent := MemberJoinedEvent{
		EventModel: midgard.EventModel{
			ID: member.GetId(),
		},
		MemberId: member.MemberId,
	}

	p.On(&memberJoinedEvent)

	return &memberJoinedEvent, nil
}

func (p *Parliament) RemoveMember(memberID MemberId) (*MemberRemovedEvent, error) {
	index := p.findIndexOfMember(memberID.ToString())

	if index == -1 {
		return nil, nil
	}

	memberRemovedEvent := MemberRemovedEvent{
		EventModel: midgard.EventModel{
			ID: memberID.ToString(),
		},
		MemberId: memberID,
	}

	p.On(&memberRemovedEvent)

	return &memberRemovedEvent, nil
}

func (p *Parliament) ValidateRepresentative(representatives []*Representative) bool {
	for _, representatives := range representatives {
		index := p.findIndexOfMember(representatives.GetIdString())

		if index == -1 {
			return false
		}
	}

	return true
}

func (p *Parliament) findIndexOfMember(memberID string) int {
	for i, member := range p.Members {
		if member.MemberId.Id == memberID {
			return i
		}
	}

	return -1
}

func (p *Parliament) FindByPeerID(memberID string) *Member {
	index := p.findIndexOfMember(memberID)

	if index == -1 {
		return nil
	}

	return p.Members[index]
}

func (p *Parliament) On(event midgard.Event) error {
	switch v := event.(type) {

	case *LeaderChangedEvent:
		p.Leader = &Leader{LeaderId:v.LeaderId}

	case *MemberJoinedEvent:
		p.Members = append(p.Members, &Member{MemberId: v.MemberId})

	case *MemberRemovedEvent:
		index := p.findIndexOfMember(v.MemberId.ToString())

		if index != -1 {
			p.Members = append(p.Members[:index], p.Members[index+1:]...)
		}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

type ParliamentRepository interface {
	GetParliament() *Parliament
	SetParliament(parliament Parliament)
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

func (pr *ParliamentRepositoryImpl) GetParliament() *Parliament {
	return &pr.parliament
}

func (pr *ParliamentRepositoryImpl) SetParliament(parliament Parliament) {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	pr.parliament = parliament
}
