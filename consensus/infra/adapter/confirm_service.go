package adapter

import (
	"errors"

	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
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

	command := consensus.CreateBlockCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Block: struct {
			Seal []byte
			Body []byte
		}{Seal: block.Seal, Body: block.Body},
	}

	return cs.publish("Command", "block.create", command)
}
