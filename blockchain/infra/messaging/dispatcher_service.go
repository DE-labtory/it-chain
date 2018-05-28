package messaging

import "github.com/it-chain/yggdrasill/impl"

type MessageDispatcher struct {}

// TODO: 추가된 block으로 BlockAddCommand 정의해야함
func (m MessageDispatcher) AddedBlock(block impl.DefaultBlock) error {
	return nil
}
