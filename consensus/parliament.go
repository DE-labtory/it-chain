package consensus

import (
	"errors"
	"fmt"

	"github.com/it-chain/midgard"
)

type LeaderId struct {
	Id string
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

type Leader struct {
	LeaderId LeaderId
}

func (l *Leader) StringLeaderId() string {
	return l.LeaderId.ToString()
}

func (l Leader) GetID() string {
	return l.StringLeaderId()
}

type MemberId struct {
	Id string
}

func (mid MemberId) ToString() string {
	return string(mid.Id)
}

type Member struct {
	MemberId MemberId
}

func (m *Member) StringMemberId() string {
	return m.MemberId.ToString()
}

func (m Member) GetID() string {
	return m.StringMemberId()
}

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

	if numOfMember >= 4 {
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
		LeaderId: leader.GetID(),
	}

	p.On(&leaderChangedEvent)

	return &leaderChangedEvent, nil
}

func (p *Parliament) AddMember(member *Member) (*MemberJoinedEvent, error) {
	if member == nil {
		return nil, errors.New("Member is nil")
	}

	if member.GetID() == "" {
		return nil, errors.New(fmt.Sprintf("Need Valid PeerID [%s]", member.GetID()))
	}

	index := p.findIndexOfMember(member.GetID())

	if index != -1 {
		return nil, errors.New(fmt.Sprintf("Already exist member [%s]", member.GetID()))
	}

	memberJoinedEvent := MemberJoinedEvent{
		EventModel: midgard.EventModel{
			ID: p.GetID(),
		},
		MemberId: member.GetID(),
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
			ID: p.GetID(),
		},
		MemberId: memberID.ToString(),
	}

	p.On(&memberRemovedEvent)

	return &memberRemovedEvent, nil
}

func (p *Parliament) ValidateRepresentative(representatives []*Representative) bool {
	for _, representatives := range representatives {
		index := p.findIndexOfMember(representatives.GetID())

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
		p.Leader = &Leader{
			LeaderId: LeaderId{v.LeaderId},
		}

	case *MemberJoinedEvent:
		p.Members = append(p.Members, &Member{
			MemberId: MemberId{v.MemberId},
		})

	case *MemberRemovedEvent:
		index := p.findIndexOfMember(v.MemberId)

		if index != -1 {
			p.Members = append(p.Members[:index], p.Members[index+1:]...)
		}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}
