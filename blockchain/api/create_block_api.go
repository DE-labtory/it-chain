package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
)

type CreateBlockApi struct {
	blockQueryApi blockchain.BlockQueryApi
	blockService  blockchain.BlockService
	publisherId   string
}

func NewCreateBlockApi(blockQueryApi blockchain.BlockQueryApi, blockService blockchain.BlockService, publisherId string) *CreateBlockApi {
	return &CreateBlockApi{
		blockQueryApi: blockQueryApi,
		blockService:  blockService,
		publisherId:   publisherId,
	}
}

func (b *CreateBlockApi) CreateBlock(txList []blockchain.Transaction) error {
	lastBlock, err := b.blockQueryApi.GetLastCommitedBlock()
	if err != nil {
		return ErrGetLastBlock
	}

	prevSeal := lastBlock.GetSeal()
	height := lastBlock.GetHeight() + 1
	defaultTxList := blockchain.ConvertTxTypeToDefaultTransaction(txList)
	creator := []byte(b.publisherId)

	block, err := blockchain.CreateProposedBlock(prevSeal, height, defaultTxList, creator)
	if err != nil {
		return err
	}

	err = b.blockService.ExecuteBlock(block)
	if err != nil {
		return err
	}

	return nil
}
