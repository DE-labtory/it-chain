package consensus

import (
	"errors"
	"fmt"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
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

func (p *Parliament) ChangeLeader(leader *Leader) error {
	if leader == nil {
		return errors.New("Leader is nil")
	}

	leaderChangedEvent := LeaderChangedEvent{
		EventModel: midgard.EventModel{
			ID: p.GetID(),
		},
		LeaderId: leader.GetID(),
	}

	if err := p.On(&leaderChangedEvent); err != nil {
		return err
	}

	if err := eventstore.Save(p.GetID(), leaderChangedEvent); err != nil {
		return err
	}

	return nil
}

func (p *Parliament) AddMember(member *Member) error {
	if member == nil {
		return errors.New("Member is nil")
	}

	if member.GetID() == "" {
		return errors.New(fmt.Sprintf("Need Valid PeerID [%s]", member.GetID()))
	}

	index := p.findIndexOfMember(member.GetID())

	if index != -1 {
		return errors.New(fmt.Sprintf("Already exist member [%s]", member.GetID()))
	}

	memberJoinedEvent := MemberJoinedEvent{
		EventModel: midgard.EventModel{
			ID: p.GetID(),
		},
		MemberId: member.GetID(),
	}

	if err := p.On(&memberJoinedEvent); err != nil {
		return err
	}

	if err := eventstore.Save(p.GetID(), memberJoinedEvent); err != nil {
		return err
	}

	return nil
}

func (p *Parliament) RemoveMember(memberID MemberId) error {
	index := p.findIndexOfMember(memberID.ToString())

	if index == -1 {
		return nil
	}

	memberRemovedEvent := MemberRemovedEvent{
		EventModel: midgard.EventModel{
			ID: p.GetID(),
		},
		MemberId: memberID.ToString(),
	}

	if err := p.On(&memberRemovedEvent); err != nil {
		return err
	}

	if err := eventstore.Save(p.GetID(), memberRemovedEvent); err != nil {
		return err
	}

	return nil
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
