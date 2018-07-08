package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")

type LeaderApi struct {
	leaderService p2p.LeaderService
	peerApi          PeerApi
}

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

func NewLeaderApi(leaderService p2p.LeaderService) LeaderApi {

	return LeaderApi{
		leaderService: leaderService,
	}
}

//todo update leader with ip address by peer!! not leader
func (la *LeaderApi) UpdateLeaderWithAddress(ipAddress string) error {

	peers := la.peerApi.GetPLTable().PeerList

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

	myPLTable := la.peerApi.GetPLTable()

	myPeerList, _ := myPLTable.GetPeerList()

	myLeader, _ := myPLTable.GetLeader()

	if len(myPeerList) < len(oppositePeerList) {

		la.leaderService.Set(oppositeLeader)

		la.peerApi.UpdatePeerList(oppositePeerList)

	} else {

		la.leaderService.Set(myLeader)

	}
	return nil
}
