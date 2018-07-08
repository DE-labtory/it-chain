package api

import "github.com/it-chain/it-chain-Engine/p2p"

type CommunicationApi struct {
	peerService          p2p.PeerService
	communicationService p2p.CommunicationService
}

func NewCommunicationApi(peerService p2p.PeerService,
	communicationService p2p.CommunicationService) CommunicationApi{
	return CommunicationApi{
		peerService:peerService,
		communicationService:communicationService,
	}
}

func (ca *CommunicationApi) DialToUnConnectedNode(peerList []p2p.Peer) error {

	//1. find unconnected peer
	//2. dial to unconnected peer
	for _, peer := range peerList {

		//err is nil if there is matching peer
		peer, err := ca.peerService.FindById(peer.PeerId)

		//dial if no peer matching peer id
		if err != nil {
			ca.communicationService.Dial(peer.IpAddress)
		}
	}

	return nil
}
