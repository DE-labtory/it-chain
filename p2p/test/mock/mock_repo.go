package mock

import "github.com/it-chain/it-chain-Engine/p2p"

func MakeFakePeerList() []p2p.Peer {

	peerList := make([]p2p.Peer, 0)
	peerList = append(peerList, p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "1",
		},
		IpAddress: "1",
	})

	peerList = append(peerList, p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "2",
		},
		IpAddress: "2",
	})

	peerList = append(peerList, p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "3",
		},
		IpAddress: "3",
	})

	return peerList
}

func MakeFakePLTable() (p2p.PLTable){

	peerList := MakeFakePeerList()
	leader := p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:"1",
		},
	}

	return p2p.PLTable{
		Leader:leader,
		PeerList:peerList,
	}
}