package api

import "github.com/it-chain/it-chain-Engine/p2p"

type ICommunicationApi interface {
	DialToUnConnectedNode(peerTable map[string]p2p.Peer) error
	DeliverPLTable(connectionId string) error
}

type CommunicationApi struct {
	pLTableQueryService  p2p.PLTableQueryService
	communicationService p2p.ICommunicationService
}

func NewCommunicationApi(pLTableQueryService p2p.PLTableQueryService, communicationService p2p.ICommunicationService) *CommunicationApi {
	return &CommunicationApi{
		pLTableQueryService:  pLTableQueryService,
		communicationService: communicationService,
	}
}

func (ca *CommunicationApi) DialToUnConnectedNode(peerTable map[string]p2p.Peer) error {

	//1. find unconnected peer
	//2. dial to unconnected peer
	for _, peer := range peerTable {

		//err is nil if there is matching peer
		peer, err := ca.pLTableQueryService.FindPeerById(peer.PeerId)

		//dial if no peer matching peer id
		if err != nil {
			ca.communicationService.Dial(peer.IpAddress)
		}
	}

	return nil
}

//Deliver Peer leader table that consists of peerList and leader
func (ca *CommunicationApi) DeliverPLTable(connectionId string) error {

	//1. get peer table
	peerTable, _ := ca.pLTableQueryService.GetPLTable()

	//2. deliver peer table
	ca.communicationService.DeliverPLTable(connectionId, peerTable)

	return nil
}
