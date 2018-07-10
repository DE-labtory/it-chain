package p2p

type LeaderQueryService interface {

	Get() Leader
}

type PeerQueryService interface {

	FindById(peerId PeerId) (Peer, error)
	FindAll() ([]Peer, error)
	FindByAddress(ipAddress string) (Peer, error)
}

type PLTableQueryService interface{

	GetPLTable() (PLTable, error)
}