package blockchain

import (
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
)

// ToDo: 삭제 제안 - junk_sound
//type SyncUpdateCommand struct {
//	midgard.EventModel
//	sync bool
//}

// ToDo: 삭제 제안 - junk_sound
//type NodeUpdateCommand struct {
//	midgard.EventModel
//	Peer
//}

type ProposeBlockCommand struct {
	midgard.CommandModel
	// TODO: Transaction이 너무 다름.
	Transactions []txpool.Transaction
}

// consensus에서 합의된 블록이 넘어오면 block pool에 저장한다.
type ConfirmBlockCommand struct {
	midgard.CommandModel
	Block Block
}

// ToDo: 삭제 제안 - junk_sound
//type BlockValidateCommand struct {
//	midgard.CommandModel
//	Block Block
//}

type GrpcDeliverCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

type GrpcReceiveCommand struct {
	midgard.CommandModel
	Body         []byte
	ConnectionID string
	Protocol     string
	FromPeer     Peer
}
