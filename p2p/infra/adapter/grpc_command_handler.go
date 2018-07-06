package adapter

import (
	"encoding/json"

	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
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
	GetPeerLeaderTable() p2p.PeerLeaderTable
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
	leaderApi       LeaderApi
	peerApi         GrpcCommandHandlerPeerApi
	peerService     GrpcCommandHandlerPeerService
	electionService p2p.ElectionService
}

func NewGrpcCommandHandler(leaderApi LeaderApi, peerApi GrpcCommandHandlerPeerApi, peerService GrpcCommandHandlerPeerService) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		leaderApi:   leaderApi,
		peerApi:     peerApi,
		peerService: peerService,
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
		_, oppositeLeader, oppositePeerList, _ := ReceivePeerLeaderTable(command.Body)

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

		peers := gch.peerApi.GetPeerList()

		for _, peer := range peers {

			if peer.IpAddress == toBeLeader.IpAddress {

				leader := p2p.Leader{
					LeaderId: p2p.LeaderId{Id: peer.PeerId.Id},
				}

				gch.leaderApi.UpdateLeader(leader)
			}

		}
	}

	return nil
}

func ReceivePeerLeaderTable(body []byte) (p2p.PeerLeaderTable, p2p.Leader, []p2p.Peer, error) {

	peerTable := p2p.PeerLeaderTable{}
	if err := json.Unmarshal(body, &peerTable); err != nil {
		//todo error 처리
		return p2p.PeerLeaderTable{}, p2p.Leader{}, []p2p.Peer{}, ErrUnmarshal
	}
	peerList, _ := peerTable.GetPeerList()
	leader, _ := peerTable.GetLeader()

	return peerTable, leader, peerList, nil
}

func UpdateWithLongerPeerList(gch *GrpcCommandHandler, oppositeLeader p2p.Leader, oppositePeerList []p2p.Peer) error {

	myPeerLeaderTable := gch.peerApi.GetPeerLeaderTable()

	myPeerList, _ := myPeerLeaderTable.GetPeerList()

	myLeader, _ := myPeerLeaderTable.GetLeader()

	if len(myPeerList) < len(oppositePeerList) {

		gch.leaderApi.UpdateLeader(oppositeLeader)

		gch.peerApi.UpdatePeerList(oppositePeerList)

	} else {

		gch.leaderApi.UpdateLeader(myLeader)

	}
	return nil
}

func DialToUnConnectedNode(peerService GrpcCommandHandlerPeerService, peerApi GrpcCommandHandlerPeerApi, peerList []p2p.Peer) error {

	for _, peer := range peerList {

		//err is nil if there is matching peer
		peer, err := peerApi.FindById(peer.PeerId)

		//dial if no peer matching peer id
		if err != nil {
			peerService.Dial(peer.IpAddress)
		}
	}

	return nil
}
