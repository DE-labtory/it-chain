package parliament

import (
	"errors"

	"fmt"

	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
)

//자기 자신을 제외한 네트워크의 Leader와 Member들
type Parliament struct {
	Leader  *Leader
	Members []*Member
}

func NewParliament() *Parliament {
	return &Parliament{
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

func (p *Parliament) SetLeader(leader *Leader) {
	p.Leader = leader
}

func (p *Parliament) AddMember(member *Member) error {

	if member.ID.ID == "" {
		return errors.New(fmt.Sprintf("Need Valid PeerID [%s]", member.ID.ID))
	}

	p.Members = append(p.Members, member)

	return nil
}

func (p *Parliament) RemoveMember(memberID PeerID) {

	index := p.findIndexOfMember(memberID)

	if index == -1 {
		return
	}

	p.Members = append(p.Members[:index], p.Members[index+1:]...)
}

//representative가 모두 Paliament에 속해있어야 한다.
func (p *Parliament) ValidateRepresentative(representatives []*consensus.Representative) bool {

	for _, representatives := range representatives {
		index := p.findIndexOfMember(PeerID{representatives.GetIdString()})

		if index == -1 {
			return false
		}
	}

	return true
}

func (p *Parliament) findIndexOfMember(memberID PeerID) int {

	for i, member := range p.Members {
		if member.ID == memberID {
			return i
		}
	}

	return -1
}
