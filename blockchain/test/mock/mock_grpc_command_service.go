package mock

import "github.com/it-chain/engine/blockchain"

type SyncCheckGrpcCommandService struct {
	SyncCheckResponseFunc func(block blockchain.Block) error
	ResponseBlockFunc     func(peerId blockchain.PeerId, block blockchain.Block) error
}

func (cs SyncCheckGrpcCommandService) SyncCheckResponse(block blockchain.Block) error {
	return cs.SyncCheckResponseFunc(block)
}
func (cs SyncCheckGrpcCommandService) ResponseBlock(peerId blockchain.PeerId, block blockchain.Block) error {
	return cs.ResponseBlockFunc(peerId, block)
}
