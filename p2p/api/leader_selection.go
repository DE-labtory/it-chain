package api

// p2p 내의 api 를 여러개 두고자 api 폴더를 만든것 같은데 그냥 peer_api 하나만 두고 사용하는건 어떨까요?
// 본 파일 명을 peer_api로 하자고 제안합니다.
// 다른 파일에서 import 하여 사용하는 경우 p2p 관련 api들을 사용할때는 어차피 /api까지만 임포트 할 것이고
// 참조시 어떤 파일에 어떤 기능이 있는것을 찾는것이 개발자로 하여금 혼동이 오는 것 같습니다.
// ---
// by frontalnh(namhoon)

import (
	"time"

	"log"

	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

type LeaderSelection struct {
	messageDispatcher           p2p.MessageDispatcher
	eventRepository  			*midgard.Repository
	nodeRepository   			NodeRepository // leveldb에 접근이 가능하도록 nodeRepository를 가짐
	leaderRepository 			p2p.LeaderRepository //leveldb의 leader 정보에 접근하기 위한 leaderRepository를 가짐
	myInfo           			*p2p.Node // 내 node 정보를 가짐.
}

func NewLeaderSelectionApi(eventRepository *midgard.Repository, repo p2p.NodeRepository, leaderRepository p2p.LeaderRepository, messageDispatcher p2p.MessageDispatcher, myInfo *p2p.Node) (*LeaderSelection, error) {
	leaderSelectionApi := &LeaderSelection{
		messageDispatcher:           messageDispatcher,
		nodeRepository:   repo,
		leaderRepository: leaderRepository,
		myInfo:           myInfo,
		eventRepository:  eventRepository,
	}

	return leaderSelectionApi, nil
}

func (ls *LeaderSelection) RequestChangeLeader() error {
	panic("implement please")
}

// todo dispatcher service 로 publish 옮기기
func (ls *LeaderSelection) RequestLeaderInfoTo(node p2p.Node) error {
	requestBody := p2p.TableRequestMessage{
		TimeUnix: time.Now().Unix(),
	}
	requestBodyByte, _ := common.Serialize(requestBody)
	deliverCommand := gateway.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{},
		Recipients:   make([]string, 0),
		Body:         requestBodyByte,
		Protocol:     "MessageDeliverCommand",
	}
	deliverCommand.Recipients = append(deliverCommand.Recipients, node.Id.ToString())
	return ls.messageDispatcher.Publisher("Command", "Messasge", deliverCommand)
}

// 리더를 바꾸기 위한 api
func (ls *LeaderSelection) changeLeader(peer *p2p.Node) error {
	events := make([]midgard.Event, 0)
	if peer.GetID() == "" {
		log.Println("need id")
		return errors.New("need Id")
	}

	events = append(events, p2p.LeaderChangeEvent{
		EventModel: midgard.EventModel{
			ID:   peer.GetID(),
			Type: "Leader",
		},
	})
	err := ls.eventRepository.Save(peer.GetID(), events...)

	if err != nil {
		return err
		log.Println(err.Error())
	}

	// todo repo(levelDB) 반영 오류시 event revert 고려 해야할 것 같음.
	ls.leaderRepository.SetLeader(p2p.Leader{
		LeaderId: p2p.LeaderId{
			Id: peer.GetID(),
		},
	})

	return nil
}
