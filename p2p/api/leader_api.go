package api

import (
	"errors"
	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")

type LeaderApi struct {
	leaderRepository   ReadOnlyLeaderRepository
	peerApi PeerApi
	myInfo             *p2p.Peer
}

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

type ReadOnlyLeaderRepository interface {

	GetLeader() p2p.Leader
}

func NewLeaderApi(
	leaderRepository ReadOnlyLeaderRepository, myInfo *p2p.Peer) *LeaderApi {

	return &LeaderApi{
		leaderRepository:   leaderRepository,
		myInfo:             myInfo,
	}
}

//todo update leader with ip address by peer!! not leader
func (leaderApi *LeaderApi) UpdateLeaderWithAddress(ipAddress string) error {

	peers := leaderApi.peerApi.GetPLTable().PeerList

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

func (leaderApi *LeaderApi) UpdateLeaderWithLongerPeerList(oppositeLeader p2p.Leader, oppositePeerList []p2p.Peer) error {

	myPLTable := leaderApi.peerApi.GetPLTable()

	myPeerList, _ := myPLTable.GetPeerList()

	myLeader, _ := myPLTable.GetLeader()

	if len(myPeerList) < len(oppositePeerList) {

		leaderApi.UpdateLeader(oppositeLeader)

		leaderApi.peerApi.UpdatePeerList(oppositePeerList)

	} else {

		leaderApi.UpdateLeader(myLeader)

	}
	return nil
}