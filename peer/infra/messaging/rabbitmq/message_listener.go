package rabbitmq

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
	"github.com/it-chain/it-chain-Engine/peer/api"
	"github.com/it-chain/it-chain-Engine/peer/domain/model"
	"github.com/it-chain/it-chain-Engine/peer/domain/repository"
	"github.com/it-chain/it-chain-Engine/peer/domain/service"
	"github.com/streadway/amqp"
)

type MessageListener struct {
	leaderSelectionApi *api.LeaderSelection
	peerRepository     *repository.Peer
	peerTable          *service.PeerTable
}

func NewMessageListener(leaderSelectionApi *api.LeaderSelection, repository *repository.Peer, table *service.PeerTable) *MessageListener {
	return &MessageListener{
		leaderSelectionApi: leaderSelectionApi,
		peerRepository:     repository,
		peerTable:          table,
	}
}

func (ml MessageListener) HandleConnCreateEvent(amqpMessage amqp.Delivery) {
	connCreateEevent := &event.ConnCreateEvent{}
	err := json.Unmarshal(amqpMessage.Body, connCreateEevent)

	if err != nil {
		// todo amqp error handle
	}
	newPeer := model.NewPeer(connCreateEevent.Address, model.PeerId(connCreateEevent.Id))
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
