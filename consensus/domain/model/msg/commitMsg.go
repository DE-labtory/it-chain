package msg

import (
	"encoding/json"
	"fmt"

	"github.com/it-chain/it-chain-Engine/consensus/domain/model/consensus"
)

type CommitMsg struct {
	ConsensusID consensus.ConsensusID
	SenderID    string
}

func (c CommitMsg) ToByte() ([]byte, error) {

	data, err := json.Marshal(c)
	if err != nil {
		panic(fmt.Sprintf("Error encoding : %s", err))
	}
	return data, nil
}
