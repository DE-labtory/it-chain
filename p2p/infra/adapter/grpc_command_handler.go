package adapter

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"errors"
	"github.com/it-chain/it-chain-Engine/gateway"
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
	GetPeerTable() (p2p.PeerTable)
	GetPeerList() []p2p.Peer
	UpdatePeerList(peerList []p2p.Peer) error
	DeliverPeerTable(connectionId string) error
	AddPeer(peer p2p.Peer)
}

type GrpcCommandHandler struct {
	leaderApi LeaderApi
	peerApi   GrpcCommandHandlerPeerApi
	commandService p2p.CommandService
}
func NewGrpcCommandHandler(leaderApi LeaderApi, peerApi GrpcCommandHandlerPeerApi, commandService p2p.CommandService) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		leaderApi: leaderApi,
		peerApi:   peerApi,
		commandService: commandService,
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


	case "PeerTableDeliverProtocol": //receive peer table
		//1. receive peer table
		_, oppositeLeader, oppositePeerList, _ := ReceiverPeerTable(command.Body)

		//2. update leader and peer list by info of node which has longer peer list
		UpdateWithLongerPeerList(gch, oppositeLeader, oppositePeerList)

		//3. dial according to peer table
		DialToUnConnectedNode(gch.commandService, gch.peerApi, oppositePeerList)

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

func ReceiverPeerTable(body []byte) (p2p.PeerTable, p2p.Leader, []p2p.Peer, error){
	peerTable := p2p.PeerTable{}
	if err := json.Unmarshal(body, &peerTable); err != nil {
		//todo error 처리
		return p2p.PeerTable{}, p2p.Leader{}, []p2p.Peer{},ErrUnmarshal
	}
	peerList, _ := peerTable.GetPeerList()
	leader, _ := peerTable.GetLeader()

	return peerTable, leader, peerList, nil
}

func UpdateWithLongerPeerList(gch *GrpcCommandHandler, oppositeLeader p2p.Leader, oppositePeerList []p2p.Peer){
	myPeerTable := gch.peerApi.GetPeerTable()
	myPeerList, _ := myPeerTable.GetPeerList()
	myLeader, _ := myPeerTable.GetLeader()

	if len(myPeerList) < len(oppositePeerList) {
		gch.leaderApi.UpdateLeader(oppositeLeader)
		gch.peerApi.UpdatePeerList(oppositePeerList)
	}else{
		gch.leaderApi.UpdateLeader(myLeader)
	}
}

func DialToUnConnectedNode(commandService p2p.CommandService, peerApi GrpcCommandHandlerPeerApi, peerList []p2p.Peer) error{
	myPeerList := peerApi.GetPeerList()

	for _, peer := range peerList{

		//
		for _, myPeer := range myPeerList{
			if myPeer == peer{
				break
			}
		}
		commandService.Dial(peer.IpAddress)
	}
}