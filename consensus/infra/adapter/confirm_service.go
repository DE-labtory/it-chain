package adapter

import (
	"errors"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/midgard"
)

type ConfirmService struct {
	publish Publish
}

func NewConfirmService(publish Publish) *ConfirmService {
	return &ConfirmService{
		publish: publish,
	}
}

func (cs *ConfirmService) ConfirmBlock(block consensus.ProposedBlock) error {
	if block.Seal == nil {
		return errors.New("Block hash is nil")
	}

	if block.Body == nil {
		return errors.New("There is no block")
	}

	cmd := command.ConfirmBlock{
		CommandModel: midgard.CommandModel{},
		Seal:         nil,
		Body:         nil,
	}

	return cs.publish("Command", "block.create", cmd)
}
