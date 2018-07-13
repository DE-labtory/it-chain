package blockchain

type SyncService struct {
	blockService BlockService
	peerService  PeerService
}

// TODO: 임의의 노드와 동기화 되었는지 확인한다.
func (ss *SyncService) SyncedCheck(peer Peer) SyncedState {

	// 싱크 프로세싱 스테이트 변경
	syncState := NewBlockSyncState()
	syncState.SetProgress(PROGRESSING)

	// lastBlock Get
	lastBlock, err := ss.blockService.GetLastBlock()
	if err != nil {
	}

	if lastBlock == nil {
		return UNSYNCED
	}

	// standardBlock Get
	standardBlock, err := ss.peerService.GetLastBlock(peer)
	if err != nil {
	}

	// Compare lastBlock vs standardBlock
	if lastBlock.GetHeight() < standardBlock.GetHeight() {
		return UNSYNCED
	}
	return SYNCED
}

// ToDo: 임의의 노드와 동기화하기 위해 구축 작업(동기화와 같은 의미)을 진행한다.
func (ss *SyncService) Construct(peer Peer) {

	lastBlock, err := ss.blockService.GetLastBlock()
	if err != nil {
	}

	//여기에서 이슈가 많을듯.
	if lastBlock == nil {
		retrievedBlock, err := ss.peerService.GetBlockByHeight(peer, 0)
		if err != nil {
		}
		CommitBlock(retrievedBlock)
		ss.Construct(peer)
	}

	lastHeight := lastBlock.GetHeight()

	standardBlock, err := ss.peerService.GetLastBlock(peer)
	if err != nil {
	}

	standardHeight := standardBlock.GetHeight()

	for lastHeight < standardHeight {

		targetHeight := lastHeight + 1

		retrievedBlock, err := ss.peerService.GetBlockByHeight(peer, targetHeight)
		if err != nil {
		}

		CommitBlock(retrievedBlock)

		raiseHeight(&lastHeight)

	}
}

func raiseHeight(height *BlockHeight) {
	*height++
}
