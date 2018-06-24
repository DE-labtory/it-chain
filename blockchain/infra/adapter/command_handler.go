package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/txpool"
)

type BlockchainCommandHandler struct {
	nodeApi api.NodeApi
	blockApi api.BlockApi
}

func NewBlockchainCommandHandler(blockApi api.BlockApi, nodeApi api.NodeApi) *BlockchainCommandHandler {
	return &BlockchainCommandHandler{
		nodeApi: nodeApi,
		blockApi: blockApi,
	}
}

// todo
func (b *BlockchainCommandHandler) HandleUpdateNodesCommand(command blockchain.NodeUpdateCommand) {
	panic("implement me")
}


func (handler *BlockchainCommandHandler) HandleProposeBlockCommand(cmd blockchain.ProposeBlockCommand, dispatcher blockchain.MessageDispatcher) {
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

	dispatcher.SendBlockValidateCommand(block)
}

// TODO: yggdrasill/impl/Transaction과 txpool/Transaction이 다름.
func convertTxList(txList []txpool.Transaction) ([]blockchain.Transaction, error) {
	return nil, nil
}

