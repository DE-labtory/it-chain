package rabbitmq

import (
	"time"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
	"github.com/it-chain/it-chain-Engine/peer/domain/model"
	"github.com/it-chain/it-chain-Engine/peer/domain/service"
)

type MessageProducer struct {
	Publish service.Publish
}

func NewMessageProducer(publish service.Publish) *MessageProducer {
	return &MessageProducer{
		Publish: publish,
	}
}


// 새로운 리더 정보를 받아오는 메서드이다.
func (mp *MessageProducer) RequestLeaderInfo(peer model.Peer) error {
	requestBody := event.LeaderInfoRequestCmd{
		TimeUnix: time.Now().Unix(),
	}
	requestBodyByte, _ := common.Serialize(requestBody)

	deliverEvent := &event.MessageDeliverEvent{
		Recipients: make([]string, 0),
		Body:       requestBodyByte,
		Protocol:   event.LeaderInfoRequestProtocol,
	}
	deliverEvent.Recipients = append(deliverEvent.Recipients, peer.Id.ToString())

	deliverSerialize, _ := common.Serialize(deliverEvent)
	return mp.Publish(topic.MessageDeliverEvent.String(), deliverSerialize)
}


// 단일 피어에게 새로운 리더 정보를 전달하는 메서드이다.
func (mp *MessageProducer) DeliverLeaderInfo(toPeer model.Peer, leader model.Peer) error {

	// 리더 정보를 leaderInfoBody에 담아줌
	leaderInfoBody := event.LeaderInfoPublishEvent{
		LeaderId: leader.Id.ToString(),
		Address:  leader.IpAddress,
	}

	// 리더 정보 json byte 변환
	leaderInfoBodyByte, _ := common.Serialize(leaderInfoBody)

	// 메세지 전달 이벤트 구조를 담는다.
	deliverEvent := event.MessageDeliverEvent{
		Recipients: make([]string, 0),
		Body:       leaderInfoBodyByte,
		Protocol:   event.LeaderInfoDeliverProtocol,
	}

	// 메세지를 수신할 수신자들을 지정해 준다.
	deliverEvent.Recipients = append(deliverEvent.Recipients, toPeer.Id.ToString())

	deliverSerialize, _ := common.Serialize(deliverEvent)

	// topic 과 serilized data를 받아 publish 한다.
	return mp.Publish(topic.MessageDeliverEvent.String(), deliverSerialize)
}

// 새로운 리더를 업데이트하는 메서드이다.
func (mp *MessageProducer) LeaderUpdateEvent(leader model.Peer) error {
	panic("implement me")
}
