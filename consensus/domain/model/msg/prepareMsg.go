package msg

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
)

type PrepareMsg struct {
	ConsensusID consensus.ConsensusID
	Block       consensus.Block
	SenderID    string
}

func (p PrepareMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}
