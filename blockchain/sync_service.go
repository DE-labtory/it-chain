package blockchain

type SyncService struct {
	blockQueryService BlockQueryService
}

// Check if Synchronizing blockchain with given peer is needed
func (ss *SyncService) syncedCheck(peer Peer) isSynced {

	// If nil peer is given(when i'm the first node of p2p network) : Synced
	if peer.IpAddress == "" {
		return SYNCED
	}

	// Get last block of my blockChain
	lastBlock, err := ss.blockQueryService.GetLastBlock()
	if err != nil {
	}

	// Get last block of other peer's blockChain
	standardBlock, err := ss.blockQueryService.GetLastBlockFromPeer(peer)
	if err != nil {
	}

	// Compare last block vs standard block
	if lastBlock.GetHeight() < standardBlock.GetHeight() {
		return UNSYNCED
	}

	return SYNCED
}

// Construct blockchain to synchronize with a given peer
func (ss *SyncService) construct(peer Peer) error {

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

		CommitBlock(retrievedBlock)

		raiseHeight(&lastHeight)
	}
	return nil
}

func setTargetHeight(lastHeight BlockHeight) BlockHeight {
	return lastHeight + 1
}

func raiseHeight(height *BlockHeight) {
	*height++
}

// Check if Synchronizing whth given peer is need and Construct to synchronize
func (ss *SyncService) SyncWithPeer(peer Peer) error {

	if state := ss.syncedCheck(peer); state == SYNCED {
		return nil
	}

	if err := ss.construct(peer); err != nil {
		return err
	}

	return nil
}
