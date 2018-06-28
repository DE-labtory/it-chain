package blockchain

type PeerId struct {
	Id string
}

type Peer struct {
	IpAddress string
	PeerId PeerId
}
