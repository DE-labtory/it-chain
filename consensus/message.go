package consensus

import "github.com/it-chain/it-chain-Engine/blockchain"

type PrePrepareMsg struct {
	ConsensusId string
	SenderId string
	ProposedBlock blockchain.Block
}

type PrepareMsg struct {
	ConsensusId string
	SenderId string
}

type CommitMsg struct {
	ConsensusId string
	SenderId string
}