package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/it-chain/it-chain-Engine/blockchain"
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

