package blockchain

type PeerService interface {
	PeerInnerService
	PeerOuterService
}

type PeerInnerService interface {
	GetRandomPeer() Peer
}

type PeerOuterService interface {
	GetLastBlock(peer Peer) Block
	GetBlockByHeight(peer Peer, height BlockHeight) Block
}

type BlockService interface {
	GetLastBlock() Block
}
