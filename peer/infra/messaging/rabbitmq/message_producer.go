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

func (mp *MessageProducer) DeliverLeaderInfo(toPeer model.Peer, leader model.Peer) error {
	leaderInfoBody := event.LeaderInfoPublishEvent{
		LeaderId: leader.Id.ToString(),
		Address:  leader.IpAddress,
	}
	leaderInfoBodyByte, _ := common.Serialize(leaderInfoBody)

	deliverEvent := event.MessageDeliverEvent{
		Recipients: make([]string, 0),
		Body:       leaderInfoBodyByte,
		Protocol:   event.LeaderInfoDeliverProtocol,
	}
	deliverEvent.Recipients = append(deliverEvent.Recipients, toPeer.Id.ToString())

	deliverSerialize, _ := common.Serialize(deliverEvent)
	return mp.Publish(topic.MessageDeliverEvent.String(), deliverSerialize)
}

func (mp *MessageProducer) LeaderUpdateEvent(leader model.Peer) error {
	panic("implement me")
}
