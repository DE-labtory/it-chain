package rabbitmq

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
	"github.com/it-chain/it-chain-Engine/peer/api"
	"github.com/it-chain/it-chain-Engine/peer"
	"github.com/streadway/amqp"
)

// 새로운 디렉토리 구조에 맞추어 peer와 repository, service 참조를 없애었습니다.

type MessageListener struct {
	leaderSelectionApi *api.LeaderSelection
	peerRepository     *PeerRepository
	peerTable          *PeerTable
}


func NewMessageListener(leaderSelectionApi *api.LeaderSelection, repository *PeerRepository, table *PeerTable) *MessageListener {
	return &MessageListener{
		leaderSelectionApi: leaderSelectionApi,
		peerRepository:     repository,
		peerTable:          table,
	}
}

// connection이 발생하면 처리하는 메소드이다.
// connection이 발생하면 peer db에 peer 를 추가한다.
func (ml MessageListener) HandleConnCreateEvent(amqpMessage amqp.Delivery) {
	connCreateEevent := &event.ConnCreateEvent{}
	err := json.Unmarshal(amqpMessage.Body, connCreateEevent)

	if err != nil {
		// todo amqp error handle
	}
	newPeer := NewPeer(connCreateEevent.Address, PeerId(connCreateEevent.Id))
	(*ml.peerRepository).Save(*newPeer)
	if ml.peerTable.GetLeader() == nil {
		err = ml.leaderSelectionApi.RequestLeaderInfoTo(*newPeer)
		if err != nil {
			// todo amqp error handle
		}
	}
}

func (ml MessageListener) HandleMessageReceiveEvent(amqpMessage amqp.Delivery) {
	receiveEvent := &event.MessageReceiveEvent{}
	err := json.Unmarshal(amqpMessage.Body, receiveEvent)
	if err != nil {
		// todo amqp error handle
	}
	// handle 해야될거만 확인 아니면 버려~
	if receiveEvent.Protocol == topic.LeaderInfoRequestCmd.String() {
		ml.peerTable.GetLeader()
	} else if receiveEvent.Protocol == topic.LeaderInfoPublishEvent.String() {

	}

}
