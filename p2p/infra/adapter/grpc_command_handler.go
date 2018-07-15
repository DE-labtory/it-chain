package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
)

var ErrLeaderInfoDeliver = errors.New("leader info deliver failed")
var ErrPeerListDeliver = errors.New("peer list deliver failed")
var ErrPeerDeliver = errors.New("peer deliver failed")
var ErrUnmarshal = errors.New("error during unmarshal")

type LeaderApi interface {
	UpdateLeaderWithAddress(ipAddress string) error
	DeliverLeaderInfo(connectionId string)
}

type GrpcCommandHandlerPeerApi interface {
	GetPLTable() p2p.PLTable
	GetPeerList() []p2p.Peer
	FindById(peerId p2p.PeerId) (p2p.Peer, error)
	UpdatePeerList(peerList []p2p.Peer) error
	DeliverPLTable(connectionId string) error
	AddPeer(peer p2p.Peer)
	UpdateLeaderWithLongerPeerList(leader p2p.Leader, peerList []p2p.Peer) error
}

type GrpcCommandHandlerCommunicationService interface {
	Dial(ipAddress string) error
}

type GrpcCommandHandler struct {
	leaderApi        LeaderApi
	peerApi          GrpcCommandHandlerPeerApi
	electionService  p2p.ElectionService
	communicationApi api.CommunicationApi
	pLTableService   p2p.PLTableService
}

func NewGrpcCommandHandler(leaderApi LeaderApi, peerApi GrpcCommandHandlerPeerApi, peerService GrpcCommandHandlerCommunicationService) *GrpcCommandHandler {
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


	case "PLTableDeliverProtocol": //receive peer table

		//1. receive peer table
		pLTable, _ := gch.pLTableService.GetPLTableFromCommand(command)

		//2. update leader and peer list by info of node which has longer peer list
		gch.peerApi.UpdateLeaderWithLongerPeerList(pLTable.Leader, pLTable.PeerList)

		//3. dial according to peer table
		gch.communicationApi.DialToUnConnectedNode(pLTable.PeerList)

		break

	case "RequestVoteProtocol":
		gch.electionService.Vote(command.ConnectionID)

	case "VoteLeaderProtocol":
		//	1. if candidate, reset left time
		//	2. count up
		//	3. if counted is same with num of peer-1 set leader and publish
		gch.electionService.DecideToBeLeader(command)

	case "UpdateLeaderProtocol":

		toBeLeader := p2p.Peer{}
		err := common.Deserialize(command.Body, toBeLeader)

		if err != nil {
			return err
		}

		gch.leaderApi.UpdateLeaderWithAddress(toBeLeader.IpAddress)
	}

	return nil
}

//
