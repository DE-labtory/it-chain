package service

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/consensus"
	"github.com/it-chain/it-chain-Engine/consensus/legacy/domain/model/parliament"
)

func Elect(parliament parliament.Parliament) ([]*consensus.Representative, error) {

	Representatives := make([]*consensus.Representative, 0)

	if parliament.GetLeader() == nil {
		return nil, errors.New("No Leader")
	}

	Representatives = append(Representatives, consensus.NewRepresentative(parliament.GetLeader().GetStringID()))

	for _, member := range parliament.Members {
		Representatives = append(Representatives, consensus.NewRepresentative(member.GetStringID()))
	}

	return Representatives, nil
}
