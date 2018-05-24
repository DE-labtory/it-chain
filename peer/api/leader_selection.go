package api

import (
	"github.com/it-chain/it-chain-Engine/peer/domain/model"
	"github.com/it-chain/it-chain-Engine/peer/domain/repository"
	"github.com/it-chain/it-chain-Engine/peer/domain/service"
)

//todo peerTableService 안에서 messageProduce 까지해줘야할듯? transaction 문제생김
//todo 예를들어, peerTableService 는 Leader 를 변경했는데 messageProduce 에 실패한경우?
//todo 이경우 rabbitMQ 에러 핸들링이라 어쩔 수 없다는 @junbeomlee 의 의견. issue 로 확인필요할듯
type LeaderSelection struct {
	peerTableService *service.PeerTable
	messageProducer  service.MessageProducer
	peerRepository   repository.Peer
	myInfo           *model.Peer
}

func NewLeaderSelectionApi(repo repository.Peer, messageProducer service.MessageProducer, myInfo *model.Peer) (*LeaderSelection, error) {
	leaderSelectionApi := &LeaderSelection{
		peerTableService: service.NewPeerTableService(repo, myInfo),
		messageProducer:  messageProducer,
		peerRepository:   repo,
		myInfo:           myInfo,
	}

	return leaderSelectionApi, nil
}

func (ls *LeaderSelection) RequestChangeLeader() error {
	panic("implement please")
}

func (ls *LeaderSelection) changeLeader(peer *model.Peer) error {
	err := ls.peerTableService.SetLeader(peer)
	if err != nil {
		return err
	}
	return ls.messageProducer.LeaderUpdateEvent(*peer)
}
