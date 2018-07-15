package service

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
)

type SyncService struct {
	blockQueryService blockchain.BlockQueryService
}

// Check if Synchronizing whth given peer is need and Construct to synchronize
func (ss *SyncService) SyncWithPeer(peer blockchain.Peer) error {

	state, err := ss.syncedCheck(peer)

	if err != nil {
		return blockchain.ErrSyncedCheck
	}

	if state == blockchain.SYNCED {
		return nil
	}

	if err := ss.construct(peer); err != nil {
		return blockchain.ErrConstruct
	}

	return nil
}

// Check if Synchronizing blockchain with given peer is needed
func (ss *SyncService) syncedCheck(peer blockchain.Peer) (blockchain.IsSynced, error) {

	// If nil peer is given(when i'm the first node of p2p network) : Synced
	if peer.IpAddress == "" {
		return blockchain.SYNCED, nil
	}

	// Get last block of my blockChain
	lastBlock, err := ss.blockQueryService.GetLastBlock()
	if err != nil {
		return blockchain.UNSYNCED, err
	}

	// Get last block of other peer's blockChain
	standardBlock, err := ss.blockQueryService.GetLastBlockFromPeer(peer)
	if err != nil {
		return blockchain.UNSYNCED, err
	}

	// Compare last block vs standard block
	if lastBlock.GetHeight() < standardBlock.GetHeight() {
		return blockchain.UNSYNCED, nil
	}

	return blockchain.SYNCED, nil
}

// Construct blockchain to synchronize with a given peer
func (ss *SyncService) construct(peer blockchain.Peer) error {

	// Get last block of my blockChain
	lastBlock, err := ss.blockQueryService.GetLastBlock()
	if err != nil {
	}

	// Set last height
	lastHeight := lastBlock.GetHeight()

	// Get last block of other peer's blockChain
	standardBlock, err := ss.blockQueryService.GetLastBlockFromPeer(peer)
	if err != nil {
	}

	// Set standard height
	standardHeight := standardBlock.GetHeight()

	// Get blocks from other peer's blockchain and commit them
	for lastHeight < standardHeight {

		targetHeight := setTargetHeight(lastHeight)

		retrievedBlock, err := ss.blockQueryService.GetBlockByHeightFromPeer(peer, targetHeight)
		if err != nil {
			return err
		}

		createdBlock, err := blockchain.CreateRetrievedBlock(retrievedBlock)
		if err != nil {
			return err
		}

		blockchain.CommitBlock(createdBlock)

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
