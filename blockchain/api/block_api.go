package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
)

type BlockApi struct {
	blockRepository blockchain.BlockRepository
}

func NewBlockApi(br blockchain.BlockRepository) BlockApi {
	return BlockApi{
		blockRepository: br,
	}
}

func (bApi BlockApi) AddBlock(block blockchain.Block) error {
	return nil
}
