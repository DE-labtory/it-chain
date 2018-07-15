package adapter

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type BlockExecuteService struct {
	publisher Publisher
}

func NewBlockExecuteService(publisher Publisher) BlockExecuteService {
	return BlockExecuteService{
		publisher: publisher,
	}
}

// TODO: When SOLO mode, Send BlockExeucteCommand to icode, otherwise send it to Consensus
func (c BlockExecuteService) ExecuteBlock(block blockchain.Block) error {
	if !blockchain.IsBlockHasAllProperties(block) {
		return errors.New("Block has missing properties")
	}
	data, err := json.Marshal(block)
	if err != nil {
		return errors.New("Error in marshaling block")
	}

	command := blockchain.BlockExecuteCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Block: data,
	}

	return c.publisher("Command", "block.execute", command)
}
