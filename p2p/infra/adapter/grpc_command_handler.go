package adapter

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"errors"
	"github.com/it-chain/it-chain-Engine/p2p/api"
)

var ErrLeaderInfoDeliver = errors.New("leader info deliver failed")
var ErrPeerListDeliver = errors.New("peer list deliver failed")
var ErrPeerDeliver = errors.New("peer deliver failed")
var ErrUnmarshal = errors.New("error during unmarshal")


type LeaderApi interface {
	UpdateLeader(leader p2p.Leader) error
	DeliverLeaderInfo(connectionId string)
}

type GrpcCommandHandlerPeerApi interface {
	GetPeerLeaderTable() (p2p.PeerLeaderTable)
	GetPeerList() []p2p.Peer
	FindById(peerId p2p.PeerId) (p2p.Peer, error)
	UpdatePeerList(peerList []p2p.Peer) error
	DeliverPeerLeaderTable(connectionId string) error
	AddPeer(peer p2p.Peer)
}

type GrpcCommandHandlerPeerService interface {
	Dial(ipAddress string) error
}
type GrpcCommandHandler struct {
	leaderApi LeaderApi
	peerApi   GrpcCommandHandlerPeerApi
	leaderRepository api.ReadOnlyLeaderRepository
	peerService GrpcCommandHandlerPeerService
	grpcCommandService GrpcCommandService
}
func NewGrpcCommandHandler(leaderApi LeaderApi, peerApi GrpcCommandHandlerPeerApi, leaderRepository api.ReadOnlyLeaderRepository, peerService GrpcCommandHandlerPeerService, grpcCommandService GrpcCommandService) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		leaderApi: leaderApi,
		leaderRepository:leaderRepository,
		peerApi:   peerApi,
		peerService: peerService,
		grpcCommandService:grpcCommandService,
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

	case "PeerLeaderTableDeliverProtocol": //receive peer table

		//1. receive peer table
		_, oppositeLeader, oppositePeerList, _ := ReceiverPeerLeaderTable(command.Body)

		//2. update leader and peer list by info of node which has longer peer list
		UpdateWithLongerPeerList(gch, oppositeLeader, oppositePeerList)

		//3. dial according to peer table
		DialToUnConnectedNode(gch.peerService, gch.peerApi, oppositePeerList)

		break

	case "PeerDeliverProtocol":

		peer := p2p.Peer{}
		err := common.Deserialize(command.Body, peer)

		if err != nil {
			return ErrPeerDeliver
		}

		gch.peerApi.AddPeer(peer)
		break

	case "RequestVoteProtocol":

		//	1. if leftTime >0, reset left time and send VoteLeaderMessage
		if gch.leaderRepository.GetLeftTime() > 0{

			gch.leaderRepository.ResetLeftTime()

			gch.grpcCommandService.DeliverVoteLeaderMessage(command.ConnectionID)

		}

	case "VoteLeaderProtocol":

		//	1. if candidate, reset left time
		//	2. count up
		if gch.leaderRepository.GetState() == "candidate"{

			gch.leaderRepository.ResetLeftTime()

			gch.leaderRepository.CountUp()

		}

		//	3. if counted is same with num of peer-1 set leader and publish
		numOfPeers := len(gch.peerApi.GetPeerList())

		if gch.leaderRepository.GetVoteCount() == numOfPeers-1{

			gch.grpcCommandService.

		}
	}

	return nil
}



func ReceiverPeerLeaderTable(body []byte) (p2p.PeerLeaderTable, p2p.Leader, []p2p.Peer, error){
	peerTable := p2p.PeerLeaderTable{}
	if err := json.Unmarshal(body, &peerTable); err != nil {
		//todo error 처리
		return p2p.PeerLeaderTable{}, p2p.Leader{}, []p2p.Peer{},ErrUnmarshal
	}
	peerList, _ := peerTable.GetPeerList()
	leader, _ := peerTable.GetLeader()

	return peerTable, leader, peerList, nil
}

func UpdateWithLongerPeerList(gch *GrpcCommandHandler, oppositeLeader p2p.Leader, oppositePeerList []p2p.Peer) error{
	myPeerLeaderTable := gch.peerApi.GetPeerLeaderTable()
	myPeerList, _ := myPeerLeaderTable.GetPeerList()
	myLeader, _ := myPeerLeaderTable.GetLeader()

	if len(myPeerList) < len(oppositePeerList) {
		gch.leaderApi.UpdateLeader(oppositeLeader)
		gch.peerApi.UpdatePeerList(oppositePeerList)
	}else{
		gch.leaderApi.UpdateLeader(myLeader)
	}
	return nil
}

func DialToUnConnectedNode(peerService GrpcCommandHandlerPeerService, peerApi GrpcCommandHandlerPeerApi, peerList []p2p.Peer) error{

	for _, peer := range peerList{
		//err is nil if there is matching peer
		peer, err := peerApi.FindById(peer.PeerId)

		//dial if no peer matching peer id
		if err !=nil{
			peerService.Dial(peer.IpAddress)
		}
	}

	return nil
}