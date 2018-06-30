package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/txpool"
	"errors"
	"sync"
)

var ErrBlockNil = errors.New("Block nil error");

type BlockApi interface {
	AddBlockToPool(block blockchain.Block)
}

type CommandHandler struct {
	blockApi BlockApi
	// block pool에 들어가는 block의 순서는 보장되어야한다.
	confirmBlockMux *sync.RWMutex
}

func NewCommandHandler(blockApi BlockApi) *CommandHandler {
	return &CommandHandler{
		confirmBlockMux: &sync.RWMutex{},
		blockApi: blockApi,
	}
}

// txpool에서 받은 transactions들을 block으로 만들어서 consensus에 보내준다.
func (h *CommandHandler) HandleProposeBlockCommand(command blockchain.ProposeBlockCommand) {
	//rawTxList := command.Transactions
	//
	//txList, err := convertTxList(rawTxList)
	//if err != nil {
	//	// TODO: handle errors
	//	return
	//}
	//
	//block, err := handler.blockApi.CreateBlock(txList)
	//if err != nil {
	//	// TODO: handle errors
	//	return
	//}
	// TODO: service는 api에서 호출되어야한다.
	//dispatcher.SendBlockValidateCommand(block)
}

// TODO: yggdrasill/impl/Transaction과 txpool/Transaction이 다름.
func convertTxList(txList []txpool.Transaction) ([]blockchain.Transaction, error) {
	return nil, nil
}

/// 합의된 block이 넘어오면 block pool에 저장한다.
func (h *CommandHandler) HandleConfirmBlockCommand(command blockchain.ConfirmBlockCommand) error {
	h.confirmBlockMux.Lock()

	defer h.confirmBlockMux.Unlock()

	block := command.Block
	if block == nil {
		return ErrBlockNil
	}

	h.blockApi.AddBlockToPool(block)
	return nil
}
