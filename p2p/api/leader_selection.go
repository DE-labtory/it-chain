package api

// p2p 내의 api 를 여러개 두고자 api 폴더를 만든것 같은데 그냥 peer_api 하나만 두고 사용하는건 어떨까요?
// 본 파일 명을 peer_api로 하자고 제안합니다.
// 다른 파일에서 import 하여 사용하는 경우 p2p 관련 api들을 사용할때는 어차피 /api까지만 임포트 할 것이고
// 참조시 어떤 파일에 어떤 기능이 있는것을 찾는것이 개발자로 하여금 혼동이 오는 것 같습니다.
// ---
// by frontalnh(namhoon)

import (
	"github.com/it-chain/it-chain-Engine/p2p/"
	"github.com/it-chain/it-chain-Engine/txpool/domain/service"
)

//todo peerTableService 안에서 messageProduce 까지해줘야할듯? transaction 문제생김
//todo 예를들어, peerTableService 는 Leader 를 변경했는데 messageProduce 에 실패한경우?
//todo 이경우 rabbitMQ 에러 핸들링이라 어쩔 수 없다는 @junbeomlee 의 의견. issue 로 확인필요할듯
type LeaderSelection struct {
	peerTableService *p2p.NodeTable          // 참여 피어에 대한 주관적인 정보를 담고 있는 peerTableService를 필드로 가짐
	messageProducer  service.MessageProducer // 메세지를 전달할 수 있는 messageProducer를 필드로 가짐
	peerRepository   p2p.NodeRepository      // leveldb에 접근이 가능하도록 peerRepository를 가짐
	myInfo           *p2p.Node               // 내 피어 정보를 가짐.
}

func NewLeaderSelectionApi(repo p2p.NodeRepository, messageProducer service.MessageProducer, myInfo *p2p.Node) (*LeaderSelection, error) {
	leaderSelectionApi := &LeaderSelection{
		peerTableService: p2p.NewNodeTableService(repo, myInfo),
		messageProducer:  messageProducer,
		peerRepository:   repo,
		myInfo:           myInfo,
	}

	return leaderSelectionApi, nil
}

func (ls *LeaderSelection) RequestChangeLeader() error {
	panic("implement please")
}

func (ls *LeaderSelection) RequestLeaderInfoTo(peer p2p.Node) error {
	return ls.messageProducer.RequestLeaderInfo(peer)
}

// 리더를 바꾸기 위한 api
func (ls *LeaderSelection) changeLeader(peer *p2p.Node) error {
	err := ls.peerTableService.SetLeader(peer)
	if err != nil {
		return err
	}
	return ls.messageProducer.LeaderUpdateEvent(*peer)
}
