package api

import "github.com/it-chain/it-chain-Engine/p2p"

type CommunicationApi struct {
	peerQueryService     p2p.PeerQueryService
	communicationService p2p.ICommunicationService
}

func NewCommunicationApi(peerQueryService p2p.PeerQueryService, communicationService p2p.ICommunicationService) CommunicationApi {
	return CommunicationApi{
		peerQueryService:peerQueryService,
		communicationService: communicationService,
	}
}

func (ca *CommunicationApi) DialToUnConnectedNode(peerList []p2p.Peer) error {

	//1. find unconnected peer
	//2. dial to unconnected peer
	for _, peer := range peerList {

		//err is nil if there is matching peer
		peer, err := ca.peerQueryService.FindById(peer.PeerId)

		//dial if no peer matching peer id
		if err != nil {
			ca.communicationService.Dial(peer.IpAddress)
		}
	}

	return nil
}
