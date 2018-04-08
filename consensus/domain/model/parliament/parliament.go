package parliament

import (
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

	if len(p.Members) >= 1 {
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

func (p *Parliament) AddMember(member *Member) {
	p.Members = append(p.Members, member)
}

func (p *Parliament) RemoveMember(memberID PeerID) {

	index := p.findIndexOfMember(memberID)

	if index == -1 {
		return
	}

	p.Members = append(p.Members[:index], p.Members[index+1:]...)
}

func (p *Parliament) ValidateRepresentative(representatives []*consensus.Representative) bool {
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
