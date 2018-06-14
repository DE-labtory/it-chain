package messaging

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
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

func (g *GrpcMessageHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) {

	switch command.Protocol {
	case "LeaderInfoRequestProtocol":
		g.leaderApi.DeliverLeaderInfo(command.FromNode.NodeId)
		break

	case "LeaderInfoDeliverProtocol":
		leader := p2p.Leader{}
		if err := json.Unmarshal(command.Data, &leader); err != nil {
			//todo error 처리
			return
		}

		g.leaderApi.UpdateLeader(leader)
		break

	case "NodeListRequestProtocol":

		g.nodeApi.DeliverNodeList(command.FromNode.NodeId)
		break

	case "NodeListDeliverProtocol":

		nodeList := make([]p2p.Node, 0)
		if err := json.Unmarshal(command.Data, &nodeList); err != nil {
			//todo error 처리
			return
		}

		g.nodeApi.UpdateNodeList(nodeList)
		break

	case "NodeDeliverProtocol":

		node := p2p.Node{}
		err := common.Deserialize(command.Data, node)

		if err != nil {
			return
		}

		g.nodeApi.AddNode(node)
		break
	}
}
