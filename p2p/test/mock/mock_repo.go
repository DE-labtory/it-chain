package mock

import "github.com/it-chain/engine/p2p"

func MakeFakePeerTable() map[string]p2p.Peer {

	peerTable := make(map[string]p2p.Peer)

	peerTable["1"] = p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "1",
		},
		IpAddress: "1",
	}
	peerTable["2"] = p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "2",
		},
		IpAddress: "2",
	}
	peerTable["3"] = p2p.Peer{
		PeerId: p2p.PeerId{
			Id: "3",
		},
		IpAddress: "3",
	}

	return peerTable
}

func MakeFakePLTable() p2p.PLTable {

	peerTable := MakeFakePeerTable()
	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{
			Id: "1",
		},
	}

	return p2p.PLTable{
		Leader:    leader,
		PeerTable: peerTable,
	}
}
