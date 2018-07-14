package blockchain

type PeerId struct {
	Id string
}

func (peerId PeerId) ToString() string {
	return string(peerId.Id)
}

type Peer struct {
	IpAddress string
	PeerId    PeerId
}
