package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyPeerTable = errors.New("empty peer list proposed")

type PeerApi struct {
	pLTableQueryService  p2p.PLTableQueryService
	communicationService p2p.ICommunicationService
}

func NewPeerApi(
	pLTableQueryService p2p.PLTableQueryService,
	communicationService p2p.ICommunicationService) *PeerApi {

	return &PeerApi{
		pLTableQueryService:pLTableQueryService,
		communicationService: communicationService,
	}
}



//Deliver Peer leader table that consists of peerList and leader
func (peerApi *PeerApi) DeliverPLTable(connectionId string) error {

	//1. get peer table
	peerTable, _ := peerApi.pLTableQueryService.GetPLTable()

	//2. deliver peer table
	peerApi.communicationService.DeliverPLTable(connectionId, peerTable)

	return nil
}
