package api

import (
	"github.com/it-chain/it-chain-Engine/consensus/domain/model/parliament"
	"github.com/it-chain/it-chain-Engine/consensus/domain/repository"
)

type ParliamentApi struct {
	parlimentRepository repository.ParlimentRepository
}

func (pApi ParliamentApi) ChangeLeader(leader parliament.Leader) error {

	parliament := pApi.parlimentRepository.Get()
	parliament.SetLeader(&leader)

	pApi.parlimentRepository.Save(parliament)

	return nil
}

func (pApi ParliamentApi) AddMember(member parliament.Member) error {

	parliament := pApi.parlimentRepository.Get()
	err := parliament.AddMember(&member)

	if err != nil {
		return err
	}

	pApi.parlimentRepository.Save(parliament)

	return nil
}

func (pApi ParliamentApi) RemoveMember(memberID parliament.PeerID) error {

	parliament := pApi.parlimentRepository.Get()
	parliament.RemoveMember(memberID)

	pApi.parlimentRepository.Save(parliament)

	return nil
}
