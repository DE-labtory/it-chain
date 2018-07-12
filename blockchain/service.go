package blockchain

type SyncService interface {
	RequestBlock(peerId PeerId, height uint64) error
	ResponseBlock(peerId PeerId, block Block) error
	SyncCheckResponse(block Block) error
}
