package blockchain

type SyncService struct {
	blockQueryService BlockQueryService
}

// TODO: 임의의 노드와 동기화 되었는지 확인한다.
func (ss *SyncService) syncedCheck(peer Peer) SyncedState {

	// 내가 첫 번째 노드일 경우
	if peer == nil {
		return true
	}

	// lastBlock Get
	lastBlock, err := ss.blockQueryService.GetLastBlock()
	if err != nil {
	}

	// standardBlock Get
	standardBlock, err := ss.blockQueryService.GetLastBlockFromPeer(peer)
	if err != nil {
	}

	// Compare lastBlock vs standardBlock
	if lastBlock.GetHeight() < standardBlock.GetHeight() {
		return false
	}

	return true
}

// ToDo: 임의의 노드와 동기화하기 위해 구축 작업(동기화와 같은 의미)을 진행한다.
func (ss *SyncService) construct(peer Peer) error {

	// lastBlock Get
	lastBlock, err := ss.blockQueryService.GetLastBlock()
	if err != nil {
	}

	// lastHeight Set
	lastHeight := lastBlock.GetHeight()

	// standardBlock Get
	standardBlock, err := ss.blockQueryService.GetLastBlockFromPeer(peer)
	if err != nil {
	}

	// standardHeight Set
	standardHeight := standardBlock.GetHeight()

	for lastHeight < standardHeight {

		// transaction 단위
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

func (ss *SyncService) Sync(peer Peer) error {

	if state := ss.syncedCheck(peer); state == true {
		return nil
	}

	if err := ss.construct(peer); err != nil {
		return err
	}

	return nil
}
