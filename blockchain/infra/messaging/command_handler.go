package messaging

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/it-chain/it-chain-Engine/txpool"
)

type BlockCommandHandler struct {
	blockApi api.BlockApi
}

func NewBlockCommandHandler(blockApi api.BlockApi) *BlockCommandHandler {
	return &BlockCommandHandler{blockApi: blockApi}
}

func (handler *BlockCommandHandler) HandleProposeBlockCommand(cmd blockchain.ProposeBlockCommand, dispatcher blockchain.MessageDispatcher) {
	rawTxList := cmd.Transactions

	txList, err := convertTxList(rawTxList)
	if err != nil {
		// TODO: handle errors
		return
	}

	block, err := handler.blockApi.CreateBlock(txList)
	if err != nil {
		// TODO: handle errors
		return
	}

	dispatcher.SendBlockCreatedEvent(block)
}

// TODO: yggdrasill/impl/Transaction과 txpool/Transaction이 다름.
func convertTxList(txList []txpool.Transaction) ([]blockchain.Transaction, error) {
	return nil, nil
}
