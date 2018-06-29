package blockchain

type PeerId struct {
	Id string
}

type Peer struct {
	IpAddress string
	PeerId PeerId
}

type PeerRepository interface {
	Add(peer Peer) error
	Remove(id PeerId) error
}