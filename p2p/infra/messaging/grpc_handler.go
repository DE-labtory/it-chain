package messaging

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/p2p/api"
)

type GrpcMessageHandler struct {
	nodeRepository   		p2p.NodeRepository
	leaderRepository 		p2p.LeaderRepository
	messageDispatcher       p2p.MessageDispatcher
	eventRepository 		midgard.Repository
}

func NewGrpcMessageHandler(nodeRepo *leveldb.NodeRepository, leaderRepo *leveldb.LeaderRepository, messageDispatcher *MessageDispatcher) *GrpcMessageHandler {
	return &GrpcMessageHandler{
		nodeRepository:   nodeRepo,
		leaderRepository: leaderRepo,
		messageDispatcher:       messageDispatcher,
	}
}

//todo implement
func (gmh *GrpcMessageHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) {
	panic("need to implement")
	switch {
	case command.Protocol=="UpdateLeader":
		leaderApi := api.NewLeaderApi(gmh.nodeRepository, gmh.leaderRepository, gmh.eventRepository, gmh.messageDispatcher)
		leaderApi.UpdateLeader(command)

	case command.Protocol=="NodeListDeliver":
		nodeApi := api.NewNodeApi(gmh.nodeRepository, gmh.leaderRepository, gmh.eventRepository, gmh.messageDispatcher)
		nodeApi.UpdateNodeList(command)
	}
	/*receiveEvent := &event.MessageReceiveEvent{}
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
	}*/
}

func (gmh *GrpcMessageHandler) HandlerMessageDeliver(command p2p.MessageDeliverCommand){
	panic("implement me!")
}