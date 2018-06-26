package consensus

import "github.com/it-chain/it-chain-Engine/blockchain"

type PrePrepareMsg struct {
	ConsensusId string
	SenderId string
	ProposedBlock blockchain.Block
}