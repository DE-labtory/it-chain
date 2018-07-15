package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type BlockService struct {
	ExecuteBlockFunc func(block blockchain.Block) error
}

func (b BlockService) ExecuteBlock(block blockchain.Block) error {
	return b.ExecuteBlockFunc(block)
}
