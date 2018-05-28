package rabbitmq

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/topic"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/it-chain/it-chain-Engine/p2p/domain/model"
	"github.com/it-chain/it-chain-Engine/p2p/domain/repository"
	"github.com/it-chain/it-chain-Engine/p2p/domain/service"
	"github.com/streadway/amqp"
)

type MessageListener struct {
	leaderSelectionApi *api.LeaderSelection
	peerRepository     *repository.PeerRepository
	peerTable          *service.PeerTable
	messageProducer    *service.MessageProducer
}

func NewMessageListener(leaderSelectionApi *api.LeaderSelection, repository *repository.PeerRepository, table *service.PeerTable, producer *service.MessageProducer) *MessageListener {

	return &MessageListener{
		leaderSelectionApi: leaderSelectionApi,
		peerRepository:     repository,
		peerTable:          table,
		messageProducer:    producer,
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
		curLeader := ml.peerTable.GetLeader()
		if curLeader == nil {
			curLeader = &model.Peer{
				IpAddress: "",
				Id:        "",
			}
		}
		// todo error handle
		toPeer, _ := (*ml.peerRepository).FindById(model.PeerId(receiveEvent.SenderId))
		// todo error handle
		err = (*ml.messageProducer).DeliverLeaderInfo(*toPeer, *curLeader)

	} else if receiveEvent.Protocol == topic.LeaderInfoPublishEvent.String() {
		eventBody := &event.LeaderInfoPublishEvent{}
		// todo error handle
		err = common.Deserialize(receiveEvent.Body, eventBody)
		leader := model.NewPeer(eventBody.Address, model.PeerId(eventBody.LeaderId))
		ml.peerTable.SetLeader(leader)
	}

}
