package adapter

import (
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/midgard"
)

type Publisher func(exchange string, topic string, data interface{}) (err error)

type CommandService struct {
	publisher Publisher
}

func NewCommandService(publisher Publisher) *CommandService {
	return &CommandService{
		publisher: publisher,
	}
}

func (c *CommandService) SendBlockExecuteResultCommand(results []icode.Result, blockId string) error {
	return c.publisher("Command", "blockResult", icode.BlockResultCommand{
		CommandModel: midgard.CommandModel{
			ID: blockId,
		},
		TxResults: results,
	})
}
