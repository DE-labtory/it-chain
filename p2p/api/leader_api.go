package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")
var ErrNoMatchingPeerWithIpAddress = errors.New("no matching peer with ip address")

type LeaderApi struct {
	leaderService       p2p.ILeaderService
	pLTableQueryService p2p.PLTableQueryService
}

func NewLeaderApi(leaderService p2p.ILeaderService, pLTableQueryService p2p.PLTableQueryService) LeaderApi {

	return LeaderApi{
		leaderService:       leaderService,
		pLTableQueryService: pLTableQueryService,
	}
}

func (la *LeaderApi) UpdateLeaderWithAddress(ipAddress string) error {

	//1. loop peer list and find specific address
	//2. update specific peer as leader
	pLTable, _ := la.pLTableQueryService.GetPLTable()

	peers := pLTable.PeerList

	for _, peer := range peers {

		if peer.IpAddress == ipAddress {

			p2p.UpdateLeader(peer)

			return nil
		}

	}

	return ErrNoMatchingPeerWithIpAddress
}

func (la *LeaderApi) UpdateLeaderWithLongerPeerList(oppositeLeader p2p.Leader, oppositePeerList []p2p.Peer) error {

	myPLTable, _ := la.pLTableQueryService.GetPLTable()

	myPeerList, _ := myPLTable.GetPeerList()

	myLeader, _ := myPLTable.GetLeader()

	if len(myPeerList) < len(oppositePeerList) {

		la.leaderService.Set(oppositeLeader)

	} else {

		la.leaderService.Set(myLeader)

	}
	return nil
}
