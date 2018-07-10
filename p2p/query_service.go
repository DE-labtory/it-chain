package p2p

type LeaderQueryService interface {

	Get() Leader
}

type PeerQueryService interface {

	FindById(id PeerId) (Peer, error)
	FindAll() ([]Peer, error)
	FindByAddress(ipAddress string) (Peer, error)
}