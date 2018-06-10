package messaging

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/it-chain/it-chain-Engine/gateway"
)

type GrpcMessageHandler struct {
	leaderApi api.LeaderApi
	nodeApi   api.NodeApi
}


func NewGrpcMessageHandler(leaderApi api.LeaderApi, nodeApi api.NodeApi) *GrpcMessageHandler {
	return &GrpcMessageHandler{
		leaderApi: leaderApi,
		nodeApi:   nodeApi,
	}
}

//todo implement
//<<<<<<< HEAD:p2p/infra/messaging/command_handler.go
//func (gmh *GrpcMessageHandler) HandleMessageReceive(command gateway.MessageReceiveCommand) {
//	leaderApi := api.NewLeaderApi(*gmh.nodeRepository, gmh.leaderRepository, gmh.eventRepository, gmh.messageDispatcher)
//	nodeApi := api.NewNodeApi(gmh.nodeRepository, gmh.leaderRepository, gmh.eventRepository, gmh.messageDispatcher)
//	switch {
//

//
//	// deliver node list when requested!
//	case command.Protocol=="NodeListRequestProtocol":
//		nodeList, _ := gmh.nodeRepository.FindAll()
//		gmh.messageDispatcher.DeliverNodeList(command.FromNode, nodeList)
//

func (g *GrpcMessageHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) {

	switch command.Protocol {
	case "UpdateLeader":

		leader := p2p.Node{}
		if err := json.Unmarshal(command.Data, &leader); err != nil {
			//todo error 처리
			return
		}

		g.leaderApi.UpdateLeader(leader)
	case "LeaderInfoRequestProtocol":
		leader := gmh.leaderRepository.GetLeader()
		gmh.messageDispatcher.DeliverLeaderInfo(command.FromNode, *leader)

	case "NodeListDeliver":

		nodeList := make([]p2p.Node, 0)
		if err := json.Unmarshal(command.Data, &nodeList); err != nil {
			//todo error 처리
			return
		}

		g.nodeApi.UpdateNodeList(nodeList)
	}
}

func (g *GrpcMessageHandler) HandlerMessageDeliver(command p2p.MessageDeliverCommand) {
	panic("implement me!")
}
