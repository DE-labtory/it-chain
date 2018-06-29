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
	GetPeerTable() ([]p2p.Peer, error)
	UpdatePeerTable(peerList []p2p.Peer) error
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

func (gch *GrpcCommandHandler) HandleMessageReceive(command p2p.GrpcReceiveCommand) error {

	switch command.Protocol {
	case "LeaderInfoRequestProtocol":
		gch.leaderApi.DeliverLeaderInfo(command.ConnectionID)
		break

	case "LeaderInfoDeliverProtocol":
		leader := p2p.Leader{}
		if err := json.Unmarshal(command.Body, &leader); err != nil {
			//todo error 처리
			return ErrLeaderInfoDeliver
		}

		gch.leaderApi.UpdateLeader(leader)
		break

	case "PeerListRequestProtocol":

		gch.peerApi.DeliverPeerList(command.ConnectionID)
		break


	case "PeerTableDeliverProtocol":
		//receive peer table

		//1. receive peer table
		oppositePeerTable := make([]p2p.Peer, 0)
		if err := json.Unmarshal(command.Body, &oppositePeerTable); err != nil {
			//todo error 처리
			return ErrPeerListDeliver
		}


		gch.peerApi.UpdatePeerTable(oppositePeerTable)

		//2. update leader
		myPeerTable, myErr := gch.peerApi.GetPeerTable()

		if len(myPeerTable) < len(oppositePeerTable) {

		}

		//3. dial according to peer table


		break

	case "PeerDeliverProtocol":

		peer := p2p.Peer{}
		err := common.Deserialize(command.Body, peer)

		if err != nil {
			return ErrPeerDeliver
		}

		gch.peerApi.AddPeer(peer)
		break
	}

	return nil
}
