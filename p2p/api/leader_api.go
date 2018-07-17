package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")
var ErrNoMatchingPeerWithIpAddress = errors.New("no matching peer with ip address")

type ILeaderApi interface {

	UpdateLeaderWithAddress(ipAddress string) error
	UpdateLeaderWithLargePeerTable(oppositePLTable p2p.PLTable) error
}

type LeaderApi struct {
	leaderService       p2p.ILeaderService
	peerQueryService p2p.PeerQueryService
}

func NewLeaderApi(leaderService p2p.ILeaderService, peerQueryService p2p.PeerQueryService) LeaderApi {

	return LeaderApi{
		leaderService:       leaderService,
		peerQueryService: peerQueryService,
	}
}

func (la *LeaderApi) UpdateLeaderWithAddress(ipAddress string) error {

	//1. loop peer list and find specific address
	//2. update specific peer as leader
	pLTable, _ := la.peerQueryService.GetPLTable()

	peers := pLTable.PeerTable

	for _, peer := range peers {

		if peer.IpAddress == ipAddress {

			p2p.UpdateLeader(peer)

			return nil
		}

	}

	return ErrNoMatchingPeerWithIpAddress
}

func (la *LeaderApi) UpdateLeaderWithLargePeerTable(oppositePLTable p2p.PLTable) error {

	myPLTable, _ := la.peerQueryService.GetPLTable()

	myLeader, _ := myPLTable.GetLeader()

	if len(myPLTable.PeerTable) < len(oppositePLTable.PeerTable) {

		la.leaderService.Set(oppositePLTable.Leader)

	} else {

		la.leaderService.Set(myLeader)

	}
	return nil
}
