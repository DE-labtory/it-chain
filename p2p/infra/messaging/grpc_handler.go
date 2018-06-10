package messaging

import (
	"encoding/json"

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

//todo implement
func (g *GrpcMessageHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) {

	switch command.Protocol {
	case "UpdateLeader":

		leader := p2p.Node{}
		if err := json.Unmarshal(command.Data, &leader); err != nil {
			//todo error 처리
			return
		}

		g.leaderApi.UpdateLeader(leader)

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
