package p2p

import (
	"github.com/pkg/errors"
)

var ErrEmptyLeaderId = errors.New("empty leader id")
var ErrEmptyPeerList = errors.New("empty peer list")
type PeerTable struct{
	Leader Leader
	PeerList []Peer
}

func NewPeerTable(leader Leader, peerList []Peer) *PeerTable{
	return &PeerTable{
		Leader:leader,
		PeerList:peerList,
	}
}

func(pt *PeerTable) GetLeader() (Leader, error) {
	if pt.Leader.LeaderId.Id == ""{
		return pt.Leader, ErrEmptyLeaderId
	}
	return pt.Leader, nil
}

func(pt *PeerTable) GetPeerList() ([]Peer, error){
	if len(pt.PeerList) == 0{
		return pt.PeerList, ErrEmptyPeerList
	}
	return pt.PeerList, nil
}
