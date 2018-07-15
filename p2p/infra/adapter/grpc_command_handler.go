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


type GrpcCommandHandlerCommunicationService interface {
	Dial(ipAddress string) error
}

type GrpcCommandHandler struct {
	leaderApi        api.LeaderApi
	electionService  p2p.ElectionService
	communicationApi api.CommunicationApi
	pLTableService   p2p.PLTableService
}

func NewGrpcCommandHandler(leaderApi api.LeaderApi) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		leaderApi: leaderApi,
	}
}

func (gch *GrpcCommandHandler) HandleMessageReceive(command p2p.GrpcReceiveCommand) error {

	switch command.Protocol {

	case "PLTableDeliverProtocol": //receive peer table

		//1. receive peer table
		pLTable, _ := gch.pLTableService.GetPLTableFromCommand(command)

		//2. update leader and peer list by info of node which has longer peer list
		gch.leaderApi.UpdateLeaderWithLargePeerTable(pLTable)

		//3. dial according to peer table
		gch.communicationApi.DialToUnConnectedNode(pLTable.PeerTable)

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
