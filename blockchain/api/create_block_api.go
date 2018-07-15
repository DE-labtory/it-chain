package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
)

type CreateBlockApi struct {
	blockQueryApi blockchain.BlockQueryApi
	publisherId   string
}

func NewCreateBlockApi(blockQueryApi blockchain.BlockQueryApi, publisherId string) *CreateBlockApi {
	return &CreateBlockApi{
		blockQueryApi: blockQueryApi,
		publisherId:   publisherId,
	}
}

func (b *CreateBlockApi) CreateBlock(txList []blockchain.Transaction) (blockchain.Block, error) {
	lastBlock, err := b.blockQueryApi.GetLastCommitedBlock()
	if err != nil {
		return nil, ErrGetLastBlock
	}

	prevSeal := lastBlock.GetSeal()
	height := lastBlock.GetHeight() + 1
	defaultTxList := blockchain.ConvertTxTypeToDefaultTransaction(txList)
	creator := []byte(b.publisherId)

	block, err := blockchain.CreateProposedBlock(prevSeal, height, defaultTxList, creator)
	if err != nil {
		return nil, err
	}

	return block, nil
}
