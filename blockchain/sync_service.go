package blockchain

type SyncService struct {
	blockQueryService BlockQueryService
}

// Check if Synchronizing whth given peer is need and Construct to synchronize
func (ss *SyncService) SyncWithPeer(peer Peer) error {

	state, err := ss.syncedCheck(peer)

	if err != nil {
		return ErrSyncedCheck
	}

	if state == SYNCED {
		return nil
	}

	if err := ss.construct(peer); err != nil {
		return ErrConstruct
	}

	return nil
}

// Check if Synchronizing blockchain with given peer is needed
func (ss *SyncService) syncedCheck(peer Peer) (isSynced, error) {

	// If nil peer is given(when i'm the first node of p2p network) : Synced
	if peer.IpAddress == "" {
		return SYNCED, nil
	}

	// Get last block of my blockChain
	lastBlock, err := ss.blockQueryService.GetLastBlock()
	if err != nil {
		return UNSYNCED, err
	}

	// Get last block of other peer's blockChain
	standardBlock, err := ss.blockQueryService.GetLastBlockFromPeer(peer)
	if err != nil {
		return UNSYNCED, err
	}

	// Compare last block vs standard block
	if lastBlock.GetHeight() < standardBlock.GetHeight() {
		return UNSYNCED, nil
	}

	return SYNCED, nil
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

		createdBlock, err := CreateRetrievedBlock(retrievedBlock)
		if err != nil {
			return err
		}

		CommitBlock(createdBlock)

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
