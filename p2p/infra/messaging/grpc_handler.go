package messaging

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
)

type GrpcMessageHandler struct {
	leaderApi         api.LeaderApi
	nodeApi           api.NodeApi
	messageDispatcher *MessageDispatcher
}

func NewGrpcMessageHandler(leaderApi api.LeaderApi, nodeApi api.NodeApi, messageDispatcher *MessageDispatcher) *GrpcMessageHandler {
	return &GrpcMessageHandler{
		leaderApi:         leaderApi,
		nodeApi:           nodeApi,
		messageDispatcher: messageDispatcher,
	}
}

func (g *GrpcMessageHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) {

	switch command.Protocol {
	case "LeaderInfoRequestProtocol":
		g.leaderApi.DeliverLeaderInfo(command.FromNode.NodeId)

	case "LeaderInfoDeliverProtocol":
		leader := p2p.Leader{}
		if err := json.Unmarshal(command.Data, &leader); err != nil {
			//todo error 처리
			return
		}

		g.leaderApi.UpdateLeader(leader)

	case "NodeListRequestProtocol":
		g.nodeApi.DeliverNodeList(command.FromNode.NodeId)

	case "NodeListDeliverProtocol":
		nodeList := make([]p2p.Node, 0)
		if err := json.Unmarshal(command.Data, &nodeList); err != nil {
			//todo error 처리
			return
		}

		g.nodeApi.UpdateNodeList(nodeList)

	case "NodeDeliverProtocol":
		node := p2p.Node{}
		err := common.Deserialize(command.Data, node)
		if err != nil {
			return
		}
		g.nodeApi.AddNode(node)
	}

}
