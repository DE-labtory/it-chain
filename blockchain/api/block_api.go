package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
)

type BlockApi struct {
	blockchainRepository blockchain.Repository
	eventRepository      *midgard.Repository
	publisherId          string
}

func NewBlockApi(blockchainRepository blockchain.Repository, eventRepository *midgard.Repository, publisherId string) (BlockApi, error) {
	return BlockApi{
		blockchainRepository: blockchainRepository,
		eventRepository:      eventRepository,
		publisherId:          publisherId,
	}, nil
}

// TODO: Check 과정에서 임의의 노드에게서 받은 blockchain 정보로 동기화 되었는지 확인한다.
func (bApi *BlockApi) SyncedCheck(block blockchain.Block) error {
	return nil
}
