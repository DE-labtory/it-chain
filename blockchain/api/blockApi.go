package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/repository"
	"github.com/it-chain/yggdrasill"
)

type BlockApi struct {
	yggDrasill      blockchainleveldb.YggDrasill
	blockRepository repository.BlockRepository
}

func (bApi BlockApi) CreateBlock() {

}

func (bApi BlockApi) AddBlock(block block.Block) {
	bApi.yggDrasill.AddBlock(block)
}

func (bApi BlockApi) GetLastBlock(block block.Block) {
	bApi.yggDrasill.GetLastBlock(block)
}
