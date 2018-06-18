package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type BlockService struct {
	publisher Publisher // midgard.client
}

func NewBlockService(publisher Publisher) *BlockService {
	return &BlockService{
		publisher: publisher,
	}
}

func (m BlockService) ProposeBlock(transactions []txpool.Transaction) error {

	if len(transactions) == 0 {
		return errors.New("Empty transaction list proposed")
	}

	deliverCommand := txpool.ProposeBlockCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Transactions: transactions,
	}

	return m.publisher("Command", "block.propose", deliverCommand)
}
