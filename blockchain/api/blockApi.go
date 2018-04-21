package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/repository"
)

type BlockApi struct {
	blockRepository repository.BlockRepository
}

func NewBlockApi(br repository.BlockRepository) BlockApi {
	return BlockApi{
		blockRepository: br,
	}
}

func (bApi BlockApi) CreateBlock() {

}

func (bApi BlockApi) AddBlock(block block.Block) error {
	err := bApi.blockRepository.AddBlock(block)
	return err
}

func (bApi BlockApi) GetLastBlock(block block.Block) {
	bApi.blockRepository.GetLastBlock(block)
}

func (bApi BlockApi) GetTransactionsById(id int) {
	bApi.blockRepository.GetTransactionsById(id)
}
