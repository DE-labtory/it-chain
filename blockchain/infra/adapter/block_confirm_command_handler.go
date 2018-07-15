package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
)

var ErrBlockNil = errors.New("Block nil error")

type BlockApi interface {
	StageBlock(block blockchain.Block) error
	CommitBlockFromPoolOrSync(blockId string) error
}

type CommandHandler struct {
	blockApi BlockApi
}

func NewCommandHandler(blockApi BlockApi) *CommandHandler {
	return &CommandHandler{
		blockApi: blockApi,
	}
}

/// 합의된 block이 넘어오면 block pool에 저장한다.
func (h *CommandHandler) HandleConfirmBlockCommand(command blockchain.ConfirmBlockCommand) error {
	block := command.Block
	if block == nil {
		return ErrBlockNil
	}

	h.blockApi.StageBlock(block)

	return nil
}
