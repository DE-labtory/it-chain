package p2p

import (
	"github.com/pkg/errors"
	"encoding/json"
	"github.com/it-chain/it-chain-Engine/p2p/infra/adapter"
	"sync"
)

var ErrEmptyLeaderId = errors.New("empty leader id")
var ErrEmptyPeerList = errors.New("empty peer list")

type PLTable struct {
	Leader   Leader
	PeerList []Peer
}

func NewPLTable(leader Leader, peerList []Peer) *PLTable {

	return &PLTable{
		Leader:   leader,
		PeerList: peerList,
	}
}

func (pt *PLTable) GetLeader() (Leader, error) {

	if pt.Leader.LeaderId.Id == "" {

		return pt.Leader, ErrEmptyLeaderId
	}

	return pt.Leader, nil
}

func (pt *PLTable) GetPeerList() ([]Peer, error) {

	if len(pt.PeerList) == 0 {

		return pt.PeerList, ErrEmptyPeerList
	}

	return pt.PeerList, nil
}

type PLTableService interface {

	GetPLTableFromCommand(command GrpcReceiveCommand) (PLTable, error)
	ClearPeerTable() error
}

// will be deleted after implemented in gateway api
type PLTableServiceReplica struct{
	mux sync.Mutex
	peerTable PeerTable
}

func (pLTableService *PLTableServiceReplica) GetPLTableFromCommand(command GrpcReceiveCommand) (PLTable, error) {

	peerTable := PLTable{}

	if err := json.Unmarshal(command.Body, &peerTable); err != nil {
		//todo error 처리
		return PLTable{}, adapter.ErrUnmarshal
	}

	return peerTable, nil
}

func (plts *PLTableServiceReplica) ClearPeerTable() {

	plts.mux.Lock()

	defer plts.mux.Unlock()

	for key := range plts.peerTable {

		delete(plts.peerTable, key)
	}

	plts.peerTable = make(map[string]Peer)
}
