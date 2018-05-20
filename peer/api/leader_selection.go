package api

// peer 내의 api 를 여러개 두고자 api 폴더를 만든것 같은데 그냥 peer_api 하나만 두고 사용하는건 어떨까요?
// 본 파일 명을 peer_api로 하자고 제안합니다.
// 다른 파일에서 import 하여 사용하는 경우 peer 관련 api들을 사용할때는 어차피 /api까지만 임포트 할 것이고
// 참조시 어떤 파일에 어떤 기능이 있는것을 찾는것이 개발자로 하여금 혼동이 오는 것 같습니다.
// ---
// by frontalnh(namhoon)

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

func (ls *LeaderSelection) RequestLeaderInfoTo(peer model.Peer) error {
	return ls.messageProducer.RequestLeaderInfo(peer)
}

// 리더를 바꾸기 위한 api
func (ls *LeaderSelection) changeLeader(peer *model.Peer) error {
	err := ls.peerTableService.SetLeader(peer)
	if err != nil {
		return err
	}
	return ls.messageProducer.LeaderUpdateEvent(*peer)
}
