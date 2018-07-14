package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
)

type BlockApi struct {
	syncService     blockchain.SyncService
	peerService     blockchain.PeerService
	eventRepository midgard.EventRepository
	publisherId     string
}

func NewBlockApi(syncService blockchain.SyncService, peerService blockchain.PeerService, eventRepository midgard.EventRepository, publisherId string) (BlockApi, error) {
	return BlockApi{
		syncService:     syncService,
		peerService:     peerService,
		eventRepository: eventRepository,
		publisherId:     publisherId,
	}, nil
}

// Synchronize blockchain with a peer in p2p network
func (bApi *BlockApi) Synchronize() error {

	// Set syncState : Progressing
	syncState := blockchain.NewBlockSyncState()
	syncState.SetProgress(blockchain.PROGRESSING)

	// Get a random peer in p2p network
	peer, err := bApi.peerService.GetRandomPeer()

	if err != nil {
		return err
	}

	// Synchronize blockchain with a random peer
	if err = bApi.syncService.SyncWithPeer(peer); err != nil {
		return err
	}

	// Set syncState : Done
	syncState.SetProgress(blockchain.DONE)

	return nil
}
