package blockchain

type PeerService interface {
	GetRandomPeerId() PeerId
	GetLastBlock(peerId PeerId) Block
	GetBlockByHeight(peerId PeerId, height BlockHeight) Block
}

type BlockService interface {
	GetLastBlock() Block
}
