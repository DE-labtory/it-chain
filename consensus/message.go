package consensus

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"encoding/json"
)

type PrePrepareMsg struct {
	ConsensusId string
	SenderId string
	ProposedBlock blockchain.Block
}

func (pp PrePrepareMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(pp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type PrepareMsg struct {
	ConsensusId string
	SenderId string
	ProposedBlock blockchain.Block
}

func (p PrepareMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type CommitMsg struct {
	ConsensusId string
	SenderId string
}

func (c CommitMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return data, nil
}