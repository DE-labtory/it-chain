package consensus

import (
	"errors"
	"fmt"
	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/consensus"
)

type Parliament struct {
	Leader  *Leader
	Members []*Member
}

func NewParliament() Parliament {
	return Parliament{
		Members: make([]*Member, 0),
		Leader:  nil,
	}
}

func (p Parliament) IsNeedConsensus() bool {
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

func (p Parliament) HasLeader() bool {

	if p.Leader == nil {
		return false
	}

	return true
}

func (p *Parliament) SetLeader(leader *Leader) error {

	if leader == nil {
		return errors.New("Leader is nil")
	}
	p.Leader = leader

	return nil
}

func (p *Parliament) AddMember(member *Member) error {

	if member == nil {
		return errors.New("Member is nil")
	}

	if member.GetId() == "" {
		return errors.New(fmt.Sprintf("Need Valid PeerID [%s]", member.GetId()))
	}

	index := p.findIndexOfMember(member.GetId())

	if index != -1 {
		return errors.New(fmt.Sprintf("Already exist member [%s]", member.GetId()))
	}

	p.Members = append(p.Members, member)

	return nil
}

func (p *Parliament) RemoveMember(memberID MemberId) {

	index := p.findIndexOfMember(memberID.ToString())

	if index == -1 {
		return
	}

	p.Members = append(p.Members[:index], p.Members[index+1:]...)
}

func (p *Parliament) ValidateRepresentative(representatives []*consensus.Representative) bool {

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
		if member.GetId() == memberID {
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

func (p *Parliament) GetLeader() *Leader {
	return p.Leader
}
