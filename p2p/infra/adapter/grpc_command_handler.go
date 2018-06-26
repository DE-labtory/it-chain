package adapter

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"errors"
)

var ErrLeaderInfoDeliver = errors.New("leader info deliver failed")
var ErrNodeListDeliver = errors.New("node list deliver failed")
var ErrNodeDeliver = errors.New("node deliver failed")

type LeaderApi interface {
	UpdateLeader(leader p2p.Leader) error
	DeliverLeaderInfo(nodeId p2p.NodeId)
}

type GrpcCommandHandlerNodeApi interface {
	UpdateNodeList(nodeList []p2p.Node) error
	DeliverNodeList(nodeId p2p.NodeId) error
	AddNode(node p2p.Node)
}

type GrpcCommandHandler struct {
	leaderApi LeaderApi
	nodeApi   GrpcCommandHandlerNodeApi
}
func NewGrpcCommandHandler(leaderApi LeaderApi, nodeApi GrpcCommandHandlerNodeApi) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		leaderApi: leaderApi,
		nodeApi:   nodeApi,
	}
}

func (g *GrpcCommandHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) error {

	switch command.Protocol {
	case "LeaderInfoRequestProtocol":
		g.leaderApi.DeliverLeaderInfo(command.FromNode.NodeId)
		break

	case "LeaderInfoDeliverProtocol":
		leader := p2p.Leader{}
		if err := json.Unmarshal(command.Data, &leader); err != nil {
			//todo error 처리
			return ErrLeaderInfoDeliver
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
			return ErrNodeListDeliver
		}

		g.nodeApi.UpdateNodeList(nodeList)
		break

	case "NodeDeliverProtocol":

		node := p2p.Node{}
		err := common.Deserialize(command.Data, node)

		if err != nil {
			return ErrNodeDeliver
		}
		g.nodeApi.AddNode(node)
		break
	}

	return nil
}
