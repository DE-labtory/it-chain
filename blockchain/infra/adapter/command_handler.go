package adapter

import (
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/it-chain/it-chain-Engine/blockchain"
)

type BlockchainCommandHandler struct {
	nodeApi api.NodeApi
}

func NewBlockchainCommandHandler(nodeApi api.NodeApi) *BlockchainCommandHandler {
	return &BlockchainCommandHandler{
		nodeApi: nodeApi,
	}
}

/// todo
func (b *BlockchainCommandHandler) HandleNodeUpdateCommand(command blockchain.NodeUpdateCommand) {

}