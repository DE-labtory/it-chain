package adapter

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"errors"
)

var ErrLeaderInfoDeliver = errors.New("leader info deliver failed")
var ErrPeerListDeliver = errors.New("peer list deliver failed")
var ErrPeerDeliver = errors.New("peer deliver failed")

type LeaderApi interface {
	UpdateLeader(leader p2p.Leader) error
	DeliverLeaderInfo(connectionId string)
}

type GrpcCommandHandlerPeerApi interface {
	UpdatePeerList(peerList []p2p.Peer) error
	DeliverPeerList(connectionId string) error
	AddPeer(peer p2p.Peer)
}

type GrpcCommandHandler struct {
	leaderApi LeaderApi

	peerApi   GrpcCommandHandlerPeerApi
}
func NewGrpcCommandHandler(leaderApi LeaderApi, peerApi GrpcCommandHandlerPeerApi) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		leaderApi: leaderApi,
		peerApi:   peerApi,
	}
}

func (g *GrpcCommandHandler) HandleMessageReceive(command p2p.GrpcReceiveCommand) error {

	switch command.Protocol {
	case "LeaderInfoRequestProtocol":
		g.leaderApi.DeliverLeaderInfo(command.ConnectionID)
		break

	case "LeaderInfoDeliverProtocol":
		leader := p2p.Leader{}
		if err := json.Unmarshal(command.Body, &leader); err != nil {
			//todo error 처리
			return ErrLeaderInfoDeliver
		}

		g.leaderApi.UpdateLeader(leader)
		break

	case "PeerListRequestProtocol":

		g.peerApi.DeliverPeerList(command.ConnectionID)
		break

	case "PeerListDeliverProtocol":


		peerList := make([]p2p.Peer, 0)
		if err := json.Unmarshal(command.Body, &peerList); err != nil {
			//todo error 처리
			return ErrPeerListDeliver
		}


		g.peerApi.UpdatePeerList(peerList)
		break

	case "PeerDeliverProtocol":

		peer := p2p.Peer{}
		err := common.Deserialize(command.Body, peer)

		if err != nil {
			return ErrPeerDeliver
		}

		g.peerApi.AddPeer(peer)
		break
	}

	return nil
}
