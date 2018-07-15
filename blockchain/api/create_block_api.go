package api

import "github.com/it-chain/it-chain-Engine/blockchain"

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

func (b *CreateBlockApi) CreateBlock(txList []blockchain.Transaction) error {
	lastBlock, err := b.blockQueryApi.GetLastCommitedBlock()
	if err != nil {
		return ErrGetLastBlock
	}

	//blockchain.CreateProposedBlock(lastBlock.GetSeal(), lastBlock.GetHeight() + 1, )

	return nil
}
