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

func (p *Parliament) findIndexOfMemeber(memberID string) int {
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
		index := p.findIndexOfMemeber(v.GetID())

		if index != -1{
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
	SetLeader(leader *Leader)
	GetLeader() *Leader
	AddMember(member *Member)
	RemoveMember(memberID MemberId)
	ValidateRepresentative(representatives []*Representative) bool
	findIndexOfMember(memberID string) int
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
	//numOfMember := 0

	//if p.HasLeader() {
	//	numOfMember = numOfMember + 1
	//}
	//
	//numOfMember = numOfMember + len(p.Members)
	//
	//if numOfMember >= 1 {
	//	return true
	//}

	return false
}

func (pr *ParliamentRepositoryImpl) HasLeader() bool {

	//if p.Leader == nil {
	//	return false
	//}

	return true
}

func (pr *ParliamentRepositoryImpl) SetLeader(leader *Leader) {

	//if leader == nil {
	//	return errors.New("Leader is nil")
	//}
	//p.Leader = leader

	return
}

func (pr *ParliamentRepositoryImpl) GetLeader() *Leader {
	//return p.Leader
	return &Leader{LeaderId: LeaderId{"a"}}
}

func (pr *ParliamentRepositoryImpl) AddMember(member *Member) {

	//if member == nil {
	//	return errors.New("Member is nil")
	//}
	//
	//if member.GetId() == "" {
	//	return errors.New(fmt.Sprintf("Need Valid PeerID [%s]", member.GetId()))
	//}
	//
	//index := p.findIndexOfMember(member.GetId())
	//
	//if index != -1 {
	//	return errors.New(fmt.Sprintf("Already exist member [%s]", member.GetId()))
	//}
	//
	//p.Members = append(p.Members, member)

	return
}

func (pr *ParliamentRepositoryImpl) RemoveMember(memberID MemberId) {

	//index := p.findIndexOfMember(memberID.ToString())
	//
	//if index == -1 {
	//	return
	//}
	//
	//p.Members = append(p.Members[:index], p.Members[index+1:]...)
}

func (pr *ParliamentRepositoryImpl) ValidateRepresentative(representatives []*Representative) bool {

	//for _, representatives := range representatives {
	//	index := p.findIndexOfMember(representatives.GetIdString())
	//
	//	if index == -1 {
	//		return false
	//	}
	//}

	return true
}

func (pr *ParliamentRepositoryImpl) findIndexOfMember(memberID string) int {

	//for i, member := range p.Members {
	//	if member.GetId() == memberID {
	//		return i
	//	}
	//}

	return -1
}

func (pr *ParliamentRepositoryImpl) FindByPeerID(memberID string) *Member {

	//index := p.findIndexOfMember(memberID)
	//
	//if index == -1 {
	//	return nil
	//}

	//return pr.Members[index]
	return &Member{MemberId: MemberId{"a"}}
}
