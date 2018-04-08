package msg

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
)

type PreprepareMsg struct {
	Consensus consensus.Consensus
	SenderID  string
}

func (pm PreprepareMsg) ToByte() ([]byte, error) {
	data, err := json.Marshal(pm)
	if err != nil {
		return nil, err
	}
	return data, nil
}
