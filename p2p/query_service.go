package p2p

type PLTableQueryService interface {
	GetPLTable() (PLTable, error)
	GetLeader() (Leader, error)
	FindPeerById(peerId PeerId) (Peer, error)
	FindPeerByAddress(ipAddress string) (Peer, error)
}
