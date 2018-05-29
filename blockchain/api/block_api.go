package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
)

type BlockApi struct {
	blockchainRepository blockchain.BlockRepository
	eventRepository *midgard.Repository
	publisherId string
}

func NewBlockApi(blockchainRepository blockchain.BlockRepository, eventRepository *midgard.Repository, publisherId string) BlockApi {
	return BlockApi{
		blockchainRepository: blockchainRepository,
		eventRepository: eventRepository,
		publisherId: publisherId,
	}
}

func (bApi BlockApi) AddBlock(block blockchain.Block) error {
	return nil
}
