package consensus

import "github.com/it-chain/midgard"

type CreateBlockCommand struct {
	midgard.CommandModel
}

type SendPrePrepareMsgCommand struct {
	midgard.CommandModel
	PrePrepareMsg struct {
		ConsensusId    ConsensusId
		SenderId       string
		Representative []*Representative
		ProposedBlock  ProposedBlock
	}
}

type SendPrepareMsgCommand struct {
	midgard.CommandModel
	PrepareMsg struct {
		ConsensusId ConsensusId
		SenderId    string
		BlockHash   []byte
	}
}

type SendCommitMsgCommand struct {
	midgard.CommandModel
	CommitMsg struct {
		ConsensusId ConsensusId
		SenderId    string
	}
}
