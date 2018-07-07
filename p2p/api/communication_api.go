package api

import "github.com/it-chain/it-chain-Engine/p2p"

type CommunicationApi struct{
	communicationService p2p.CommunicationService
	peerApi PeerApi
}

func (ca *CommunicationApi)DialToUnConnectedNode(peerList []p2p.Peer) error {

	for _, peer := range peerList {

		//err is nil if there is matching peer
		peer, err := ca.peerApi.FindById(peer.PeerId)

		//dial if no peer matching peer id
		if err != nil {
			ca.communicationService.Dial(peer.IpAddress)
		}
	}

	return nil
}

