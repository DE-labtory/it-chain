package adapter

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"errors"
)

var ErrLeaderInfoDeliver = errors.New("leader info deliver failed")
var ErrPeerListDeliver = errors.New("node list deliver failed")
var ErrPeerDeliver = errors.New("node deliver failed")

type LeaderApi interface {
	UpdateLeader(leader p2p.Leader) error
	DeliverLeaderInfo(nodeId p2p.PeerId)
}

type GrpcCommandHandlerPeerApi interface {
	UpdatePeerList(nodeList []p2p.Peer) error
	DeliverPeerList(nodeId p2p.PeerId) error
	AddPeer(node p2p.Peer)
}

type GrpcCommandHandler struct {
	leaderApi LeaderApi
	nodeApi   GrpcCommandHandlerPeerApi
}
func NewGrpcCommandHandler(leaderApi LeaderApi, nodeApi GrpcCommandHandlerPeerApi) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		leaderApi: leaderApi,
		nodeApi:   nodeApi,
	}
}

func (g *GrpcCommandHandler) HandleMessageReceive(command p2p.GrpcRequestCommand) error {

	switch command.Protocol {
	case "LeaderInfoRequestProtocol":
		g.leaderApi.DeliverLeaderInfo(command.FromPeer.PeerId)
		break

	case "LeaderInfoDeliverProtocol":
		leader := p2p.Leader{}
		if err := json.Unmarshal(command.Data, &leader); err != nil {
			//todo error 처리
			return ErrLeaderInfoDeliver
		}

		g.leaderApi.UpdateLeader(leader)
		break

	case "PeerListRequestProtocol":

		g.nodeApi.DeliverPeerList(command.FromPeer.PeerId)
		break

	case "PeerListDeliverProtocol":

		nodeList := make([]p2p.Peer, 0)
		if err := json.Unmarshal(command.Data, &nodeList); err != nil {
			//todo error 처리
			return ErrPeerListDeliver
		}

		g.nodeApi.UpdatePeerList(nodeList)
		break

	case "PeerDeliverProtocol":

		node := p2p.Peer{}
		err := common.Deserialize(command.Data, node)

		if err != nil {
			return ErrPeerDeliver
		}
		g.nodeApi.AddPeer(node)
		break
	}

	return nil
}
