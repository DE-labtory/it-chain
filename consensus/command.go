package consensus

import "github.com/it-chain/midgard"

type CreateBlockCommand struct {
	midgard.CommandModel
}

type SendPrePrepareMsgCommand struct {
	midgard.CommandModel
	PrePrepareMsg PrePrepareMsg
}

type SendPrepareMsgCommand struct {
	midgard.CommandModel
	PrepareMsg PrepareMsg
}

type SendCommitMsgCommand struct {
	midgard.CommandModel
	CommitMsg CommitMsg
}
