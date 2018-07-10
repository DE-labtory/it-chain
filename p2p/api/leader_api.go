package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")

type LeaderApi struct {
	leaderService  p2p.LeaderService
	pLTableService p2p.PLTableService
	pLTableApi PLTableApi
}

func NewLeaderApi(
	pLTableApi PLTableApi) LeaderApi {

	return LeaderApi{
		pLTableApi:pLTableApi,
	}
}

//todo update leader with ip address by peer!! not leader
func (la *LeaderApi) UpdateLeaderWithAddress(ipAddress string) error {

	pLTable, _ := la.pLTableApi.GetPLTable()

	peers := pLTable.PeerList

	for _, peer := range peers {

		if peer.IpAddress == ipAddress {

			leader := p2p.Leader{
				LeaderId: p2p.LeaderId{Id: peer.PeerId.Id},
			}

			p2p.UpdateLeader(leader)
		}

	}

	return nil
}

func (la *LeaderApi) UpdateLeaderWithLongerPeerList(oppositeLeader p2p.Leader, oppositePeerList []p2p.Peer) error {

	myPLTable, _ := la.pLTableApi.GetPLTable()

	myPeerList, _ := myPLTable.GetPeerList()

	myLeader, _ := myPLTable.GetLeader()

	if len(myPeerList) < len(oppositePeerList) {

		la.leaderService.Set(oppositeLeader)

	} else {

		la.leaderService.Set(myLeader)

	}
	return nil
}
