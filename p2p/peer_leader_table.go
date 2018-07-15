package p2p

import (
	"github.com/pkg/errors"
	"encoding/json"
)

var ErrEmptyLeaderId = errors.New("empty leader id")
var ErrEmptyPeerTable = errors.New("empty peer list")

type PLTable struct {
	Leader   Leader
	PeerTable map[string]Peer
}


func NewPLTable(leader Leader, peerTable map[string]Peer) *PLTable {

	return &PLTable{
		Leader:   leader,
		PeerTable: peerTable,
	}
}

func (pt *PLTable) GetLeader() (Leader, error) {

	if pt.Leader.LeaderId.Id == "" {

		return pt.Leader, ErrEmptyLeaderId
	}

	return pt.Leader, nil
}

func (pt *PLTable) GetPeerTable() (map[string]Peer, error) {

	if len(pt.PeerTable) == 0 {

		return pt.PeerTable, ErrEmptyPeerTable
	}

	return pt.PeerTable, nil
}

type PLTableService struct{}

func (plts *PLTableService) GetPLTableFromCommand(command GrpcReceiveCommand) (PLTable, error) {

	peerTable := PLTable{}

	if err := json.Unmarshal(command.Body, &peerTable); err != nil {
		//todo error 처리
		return PLTable{}, nil
	}

	return peerTable, nil
}
