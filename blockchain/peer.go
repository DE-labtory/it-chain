package blockchain

// 노드 구조체 선언.
type Peer struct {
	IpAddress string
	PeerId    PeerId
}

// PeerId 선언
type PeerId struct {
	Id string
}
