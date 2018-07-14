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

// ToDo: 임의의 노드와 블록 체인을 동기화합니다.
func (bApi *BlockApi) Synchronize() error {

	// 싱크 프로세싱 스테이트 변경
	syncState := blockchain.NewBlockSyncState()
	syncState.SetProgress(blockchain.PROGRESSING)

	// peer random으로 가져오기
	peer, err := bApi.peerService.GetRandomPeer()

	if err != nil {
		return err
	}

	// sync 하기
	if err = bApi.syncService.Sync(peer); err != nil {
		return err
	}

	syncState.SetProgress(blockchain.DONE)

	//event.save

	return nil
}
