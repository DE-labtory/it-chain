package adapter

import (
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type BlockExecuteService struct {
	publisher Publisher
}

func NewBlockExecuteService(publisher Publisher) *BlockExecuteService {
	return &BlockExecuteService{
		publisher: publisher,
	}
}

func (s *BlockExecuteService) ExecuteBlock(block blockchain.Block) error {
	command, err := createBlockExecuteCommand(block)
	if err != nil {
		return err
	}

	return s.publisher("Command", "block.execute", command)
}

func createBlockExecuteCommand(block blockchain.Block) (blockchain.BlockExecuteCommand, error){
	data, err := common.Serialize(block)
	if err != nil {
		return blockchain.BlockExecuteCommand{}, err
	}

	return blockchain.BlockExecuteCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Block: data,
	}, nil
}