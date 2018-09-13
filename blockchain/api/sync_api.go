package api

import (
	"github.com/it-chain/engine/blockchain"
)

type SyncApi struct {
	publisherId     string
	blockRepository blockchain.BlockRepository
	eventService    blockchain.EventService
	queryService    blockchain.QueryService
}

func NewSyncApi(publisherId string, blockRepository blockchain.BlockRepository, eventService blockchain.EventService, queryService blockchain.QueryService) (SyncApi, error) {
	return SyncApi{
		publisherId:     publisherId,
		blockRepository: blockRepository,
		eventService:    eventService,
		queryService:    queryService,
	}, nil
}

func (sApi SyncApi) Synchronize() error {

	// get random peer
	randomPeer, err := sApi.queryService.GetRandomPeer()

	if err != nil {
		return err
	}

	if sApi.isSynced(randomPeer) {
		return nil
	}

	// if sync has not done, on sync
	sApi.syncWithPeer(randomPeer)

	return nil

}

func (sApi SyncApi) isSynced(peer blockchain.Peer) bool {

	// If nil peer is given(when i'm the first node of p2p network) : Synced
	if peer.IpAddress == "" {
		return true
	}
	// Get last block of my blockChain
	lastBlock, err := sApi.blockRepository.FindLast()
	if err != nil {
		return false
	}
	// Get last block of other peer's blockChain
	standardBlock, err := sApi.queryService.GetLastBlockFromPeer(peer)

	if err != nil {
		return false
	}
	// Compare last block vs standard block
	if lastBlock.GetHeight() < standardBlock.GetHeight() {
		return false
	}

	return true
}

func (sApi SyncApi) syncWithPeer(peer blockchain.Peer) error {
	standardBlock, err := sApi.queryService.GetLastBlockFromPeer(peer)

	if err != nil {
		return err
	}

	standardHeight := standardBlock.GetHeight()

	lastBlock, err := sApi.blockRepository.FindLast()

	if err != nil {
		return err
	}

	lastHeight := lastBlock.GetHeight()

	return sApi.construct(peer, standardHeight, lastHeight)

}

func (sApi SyncApi) construct(peer blockchain.Peer, standardHeight blockchain.BlockHeight, lastHeight blockchain.BlockHeight) error {

	for lastHeight < standardHeight {

		targetHeight := setTargetHeight(lastHeight)
		retrievedBlock, err := sApi.queryService.GetBlockByHeightFromPeer(peer, targetHeight)

		if err != nil {
			return err
		}

		err = sApi.blockRepository.Save(retrievedBlock)

		if err != nil {
			return err
		}

		// publish
		commitEvent, err := createBlockCommittedEvent(retrievedBlock)

		if err != nil {
			return err
		}

		err = sApi.eventService.Publish("block.committed", commitEvent)
		if err != nil {
			return err
		}

		raiseHeight(&lastHeight)
	}

	return nil
}

func setTargetHeight(lastHeight blockchain.BlockHeight) blockchain.BlockHeight {
	return lastHeight + 1
}
func raiseHeight(height *blockchain.BlockHeight) {
	*height++
}
