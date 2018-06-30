package api

import (
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/consensus"
)

type ParliamentApi struct {
	eventRepository *midgard.Repository
}

func NewParliamentApi(eventRepository *midgard.Repository) ParliamentApi {
	return ParliamentApi{
		eventRepository: eventRepository,
	}
}

// todo : Implement & Event Sourcing 첨가

func (p ParliamentApi) ChangeLeader(leader consensus.Leader) error {
	return nil
}

func (p ParliamentApi) AddMember(member consensus.Member) error {
	return nil
}

func (p ParliamentApi) RemoveMember(memberId consensus.MemberId) error {
	return nil
}