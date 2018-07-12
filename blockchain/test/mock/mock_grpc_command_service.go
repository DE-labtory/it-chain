package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type MockSyncCheckGrpcCommandService struct {
	SyncCheckResponseFunc func(block blockchain.Block) error
	ResponseBlockFunc     func(peerId blockchain.PeerId, block blockchain.Block) error
}
func (cs MockSyncCheckGrpcCommandService) SyncCheckResponse(block blockchain.Block) error {
	return cs.SyncCheckResponseFunc(block)
}
func (cs MockSyncCheckGrpcCommandService) ResponseBlock(peerId blockchain.PeerId, block blockchain.Block) error {
	return cs.ResponseBlockFunc(peerId, block)
}
