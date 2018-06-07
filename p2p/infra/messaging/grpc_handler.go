package messaging

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/midgard"
	"log"
)

type GrpcMessageHandler struct {
	nodeRepository   		p2p.NodeRepository
	leaderRepository 		p2p.LeaderRepository
	messageDispatcher       *MessageDispatcher
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
	case command.Protocol=="LeaderInfoUpdate":
		if command.GetID() == ""{
			return
		}
		id := command.GetID()
		leader := &p2p.Leader{}
		err := common.Deserialize(command.Data, leader)
		if err != nil{
			panic(err)
		}

		events := make([]midgard.Event, 0)
		leaderUpdatedEvent := p2p.LeaderUpdatedEvent{
			EventModel: midgard.EventModel{
				ID: 	id,
				Type:	"Leader",
			},
			Leader: *leader,
		}

		events = append(events, leaderUpdatedEvent)
		err2 := gmh.eventRepository.Save(command.GetID(), events...)

		if err2 != nil {
			log.Println(err2.Error())
		}

		//gmh.leaderRepository.SetLeader(*leader)
		gmh.messageDispatcher.publisher.Publish("event", "leader.update", leaderUpdatedEvent)

	case command.Protocol=="NodeListDeliver":
		if command.GetID() ==""{
			return
		}

		id := command.GetID()

		nodeList := make([]p2p.Node,0)
		err := common.Deserialize(command.Data, nodeList)

		if err != nil{
			err.Error()
		}

		event := p2p.NodeListUpdatedEvent{
			EventModel: midgard.EventModel{
				ID:id,
				Type:"Node",
			},
			NodeList:nodeList,
		}

		gmh.messageDispatcher.publisher.Publish("event", "node.update", event)
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