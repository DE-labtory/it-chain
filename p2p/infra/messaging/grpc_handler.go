package messaging

import "github.com/it-chain/it-chain-Engine/p2p"

type GrpcMessageHandler struct {
	nodeRepository   p2p.NodeRepository
	leaderRepository p2p.LeaderRepository
	dispatcher       *Dispatcher
}

func NewGrpcMessageHandler(nodeRepo *p2p.NodeRepository, leaderRepo *p2p.LeaderRepository, dispatcher *Dispatcher) *GrpcMessageHandler {
	return &GrpcMessageHandler{
		nodeRepository:   nodeRepo,
		leaderRepository: leaderRepo,
		dispatcher:       dispatcher,
	}
}

//todo implement
func (gmh *GrpcMessageHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) {
	panic("need to implement")

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
