package p2p

type PeerQueryService interface {

	GetPLTable() (PLTable, error)
	GetLeader() (Leader, error)
	FindPeerById(peerId PeerId) (Peer, error)
	FindPeerByAddress(ipAddress string) (Peer, error)
}
