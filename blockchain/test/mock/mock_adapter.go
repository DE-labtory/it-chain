package mock

import "github.com/it-chain/engine/blockchain"

type BlockAdapter struct {
	GetLastBlockFunc     func(address string) (blockchain.DefaultBlock, error)
	GetBlockByHeightFunc func(address string, height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
}

func (a BlockAdapter) GetLastBlock(address string) (blockchain.DefaultBlock, error) {
	return a.GetLastBlockFunc(address)
}

func (a BlockAdapter) GetBlockByHeight(address string, height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return a.GetBlockByHeightFunc(address, height)
}
