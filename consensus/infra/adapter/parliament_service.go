package adapter

import (
	c "github.com/it-chain/it-chain-Engine/consensus"
	"errors"
)

type ParliamentService struct {

}

func NewParliamentService() *ParliamentService {
	return &ParliamentService{}
}

func (ps ParliamentService) Elect(parliament c.Parliament) ([]*c.Representative, error) {
	representatives := make([]*c.Representative, 0)

	if !parliament.HasLeader() {
		return nil, errors.New("No Leader")
	}

	representatives = append(representatives, c.NewRepresentative(parliament.Leader.GetID()))

	for _, member := range parliament.Members {
		representatives = append(representatives, c.NewRepresentative(member.GetID()))
	}

	return representatives, nil
}