package messaging

import (
	"time"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type Publisher func(exchange string, topic string, data interface{}) (err error)
type MessageDispatcher struct {
	publisher Publisher
}

func NewMessageDispatcher(publisher Publisher) *MessageDispatcher {
	return &MessageDispatcher{
		publisher: publisher,
	}
}

// 새로운 리더 정보를 받아오는 메서드이다.
func (md *MessageDispatcher) RequestLeaderInfo(peer p2p.Node) error {

	requestBody := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	requestBodyByte, _ := common.Serialize(requestBody)

	deliverCom := &p2p.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       requestBodyByte,
		Protocol:   "LeaderInfoRequestMessage",
	}
	deliverCom.Recipients = append(deliverCom.Recipients, peer.NodeId.ToString())

	return md.publisher("Command", "GrpcMessage", deliverCom)
}

// 단일 피어에게 새로운 리더 정보를 전달하는 메서드이다.
func (md *MessageDispatcher) DeliverLeaderInfo(toPeer p2p.Node, leader p2p.Node) error {

	// 리더 정보를 leaderInfoBody에 담아줌
	leaderInfoBody := p2p.LeaderInfoResponseMessage{
		LeaderId: leader.NodeId.ToString(),
		Address:  leader.IpAddress,
	}

	// 리더 정보 json byte 변환
	leaderInfoBodyByte, _ := common.Serialize(leaderInfoBody)

	// 메세지 전달 이벤트 구조를 담는다.
	deliverCommand := p2p.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       leaderInfoBodyByte,
		Protocol:   event.LeaderInfoDeliverProtocol,
	}

	// 메세지를 수신할 수신자들을 지정해 준다.
	deliverCommand.Recipients = append(deliverCommand.Recipients, toPeer.NodeId.ToString())

	// topic 과 serilized data를 받아 publisher 한다.
	return md.publisher("Command", "MessageDeliverCommand", deliverCommand)
}

func (md *MessageDispatcher) RequestTable(toNode p2p.Node) error {
	panic("implement me")
}

func (md *MessageDispatcher) ResponseTable(toNode p2p.Node, nodes []p2p.Node) error {
	panic("implement me")
}

// 새로운 리더를 업데이트하는 메서드이다.
func (md *MessageDispatcher) LeaderUpdateEvent(leader p2p.Node) error {
	panic("implement me")
}
