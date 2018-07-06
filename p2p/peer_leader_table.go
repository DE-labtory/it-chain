package p2p

import (
	"github.com/pkg/errors"
)

var ErrEmptyLeaderId = errors.New("empty leader id")
var ErrEmptyPeerList = errors.New("empty peer list")

type PeerLeaderTable struct {
	Leader   Leader
	PeerList []Peer
}

func NewPeerLeaderTable(leader Leader, peerList []Peer) *PeerLeaderTable {

	return &PeerLeaderTable{
		Leader:   leader,
		PeerList: peerList,
	}
}

func (pt *PeerLeaderTable) GetLeader() (Leader, error) {

	if pt.Leader.LeaderId.Id == "" {

		return pt.Leader, ErrEmptyLeaderId
	}

	return pt.Leader, nil
}

func (pt *PeerLeaderTable) GetPeerList() ([]Peer, error) {

	if len(pt.PeerList) == 0 {

		return pt.PeerList, ErrEmptyPeerList
	}

	return pt.PeerList, nil
}
